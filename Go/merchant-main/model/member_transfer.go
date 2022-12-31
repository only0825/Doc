package model

import (
	"errors"
	"fmt"
	"merchant/contrib/helper"
	"time"

	g "github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

// MemberTransferSubCheck 检查当前会员是否有下级
func MemberTransferSubCheck(username string) bool {

	var num int
	ex := g.Ex{
		"parent_name": username,
		"prefix":      meta.Prefix,
	}
	query, _, _ := dialect.From("tbl_members").Select(g.COUNT("uid").As("num")).Where(ex).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Get(&num, query)
	if err == nil && num == 0 {
		return false
	}

	return true
}

//MemberTransferAg 跳线转代
func MemberTransferAg(mb, destMb Member, admin map[string]string, isOfficial bool) error {

	tx, err := meta.MerchantDB.Begin() // 开启事务
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	ex := g.Ex{
		"uid":    mb.UID,
		"prefix": meta.Prefix,
	}
	record := g.Record{
		"parent_uid":  destMb.UID,
		"parent_name": destMb.Username,
		"top_uid":     destMb.TopUid,
		"top_name":    destMb.TopName,
		"tester":      destMb.Tester,
	}
	query, _, _ := dialect.Update("tbl_members").Set(record).Where(ex).ToSQL()
	fmt.Println(query)
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(err, helper.DBErr)
	}

	query = fmt.Sprintf("delete from tbl_members_tree where descendant = %s and prefix = '%s'", mb.UID, meta.Prefix)
	fmt.Println(query)
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(err, helper.DBErr)
	}

	treeNode := MemberClosureInsert(mb.UID, destMb.UID)
	fmt.Println(treeNode)
	_, err = tx.Exec(treeNode)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(err, helper.DBErr)
	}

	// 记录转代日志
	transRecord := AgencyTransferRecord{
		Id:            helper.GenLongId(),
		Flag:          551,
		Uid:           mb.UID,
		Username:      mb.Username,
		Type:          mb.AgencyType,
		BeforeUid:     mb.ParentUid,
		BeforeName:    mb.ParentName,
		AfterUid:      destMb.UID,
		AfterName:     destMb.Username,
		Remark:        "会员转代",
		UpdatedAt:     time.Now().Unix(),
		UpdatedUid:    admin["id"],
		UpdatedName:   admin["name"],
		BeforeTopUid:  mb.TopUid,
		BeforeTopName: mb.TopName,
		AfterTopUid:   destMb.TopUid,
		AfterTopName:  destMb.TopName,
		Prefix:        meta.Prefix,
	}
	query, _, _ = dialect.Insert("tbl_agency_transfer_record").Rows(transRecord).ToSQL()
	fmt.Println(query)
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	rebateRecord := g.Record{
		"parent_uid": destMb.UID,
	}
	dest := MemberRebateResult_t{}
	// 官方会员转代,默认返水比例设置为和转移到的代理一致
	if isOfficial {
		dest, err = MemberRebateFindOne(destMb.UID)
		if err != nil {
			return err
		}

		rebateRecord["zr"] = dest.ZR.StringFixed(1)
		rebateRecord["ty"] = dest.TY.StringFixed(1)
		rebateRecord["dj"] = dest.DJ.StringFixed(1)
		rebateRecord["dz"] = dest.DZ.StringFixed(1)
		rebateRecord["by"] = dest.BY.StringFixed(1)
		rebateRecord["cp"] = dest.CP.StringFixed(1)
		rebateRecord["qp"] = dest.QP.StringFixed(1)
		rebateRecord["fc"] = dest.FC.StringFixed(1)
		rebateRecord["cg_official_rebate"] = dest.CGOfficialRebate.StringFixed(2)
		rebateRecord["cg_high_rebate"] = dest.CGHighRebate.StringFixed(2)
	}
	query, _, _ = dialect.Update("tbl_member_rebate_info").Set(rebateRecord).Where(g.Ex{"uid": mb.UID}).ToSQL()
	fmt.Println(query)
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	_ = tx.Commit()

	_ = MemberRebateUpdateCache1(mb.UID, dest)

	param := map[string]interface{}{
		"uid":         mb.UID,
		"username":    mb.Username,
		"parent_uid":  destMb.UID,
		"parent_name": destMb.Username,
		"top_uid":     destMb.TopUid,
		"top_name":    destMb.TopName,
		"prefix":      meta.Prefix,
	}
	// bet 去修改注单的
	err = TransAg("transfer_ags", param)
	if err != nil {
		_ = pushLog(err, helper.ServerErr)
	}

	// task 去修改存提 调整红利的
	err = BeanPut("transfer_ag", param)
	if err != nil {
		_ = pushLog(err, helper.ServerErr)
	}

	return nil
}

//MemberTransferExist 转代申请记录是否存在
func MemberTransferExist(username string) bool {

	var num int
	ex := g.Ex{
		"username": username,
		"prefix":   meta.Prefix,
		"status":   []int{1, 2},
	}
	query, _, _ := dialect.From("tbl_agency_transfer_apply").Select(g.COUNT("uid").As("num")).Where(ex).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Get(&num, query)
	if err == nil && num == 0 {
		return false
	}

	return true
}

//MemberTransferList 团队转代申请列表
func MemberTransferList(page, pageSize int, startTime, endTime, reviewStartTime, reviewEndTime string, ex g.Ex) (AgencyTransferData, error) {

	data := AgencyTransferData{}
	// 没有查询条件  startTime endTime 必填
	if len(ex) == 0 && (startTime == "" || endTime == "") {
		return data, errors.New(helper.QueryTermsErr)
	}

	if startTime != "" && endTime != "" {

		startAt, err := helper.TimeToLoc(startTime, loc)
		if err != nil {
			return data, errors.New(helper.DateTimeErr)
		}

		endAt, err := helper.TimeToLoc(endTime, loc)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		if startAt >= endAt {
			return data, errors.New(helper.QueryTimeRangeErr)
		}

		ex["apply_at"] = g.Op{"between": exp.NewRangeVal(startAt, endAt)}
	}

	if reviewStartTime != "" && reviewEndTime != "" {

		rStart, err := helper.TimeToLoc(reviewStartTime, loc)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		rEnd, err := helper.TimeToLoc(reviewEndTime, loc)
		if err != nil {
			return data, errors.New(helper.TimeTypeErr)
		}

		if rStart >= rEnd {
			return data, errors.New(helper.QueryTimeRangeErr)
		}

		ex["review_at"] = g.Op{"between": exp.NewRangeVal(rStart, rEnd)}
	}

	t := dialect.From("tbl_agency_transfer_apply")
	if page == 1 {
		query, _, _ := t.Select(g.COUNT("id")).Where(ex).ToSQL()
		fmt.Println(query)
		err := meta.MerchantDB.Get(&data.T, query)
		if err != nil {
			return data, pushLog(err, helper.DBErr)
		}

		if data.T == 0 {
			return data, nil
		}
	}

	offset := pageSize * (page - 1)
	query, _, _ := t.Select(colsAgencyTransfer...).Where(ex).
		Offset(uint(offset)).Limit(uint(pageSize)).Order(g.C("apply_at").Desc()).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Select(&data.D, query)
	if err != nil {
		return data, pushLog(err, helper.DBErr)
	}

	return data, nil
}

//MemberTransferInsert 团队转代申请提交
func MemberTransferInsert(mb, destMb Member, admin map[string]string, remark string) error {

	record := g.Record{
		"id":            helper.GenId(),
		"prefix":        meta.Prefix,
		"uid":           mb.UID,            //转代会员uid
		"username":      mb.Username,       //转代人名
		"before_uid":    mb.ParentUid,      //转代前上级代理uid
		"before_name":   mb.ParentName,     //转代前上级代理名
		"after_uid":     destMb.UID,        //转代后上级代理uid
		"after_name":    destMb.Username,   //转代后上级代理名
		"status":        1,                 //状态 1审核中 2审核通过 3审核拒绝 4删除
		"apply_at":      time.Now().Unix(), //添加时间
		"apply_uid":     admin["id"],       //添加人uid
		"apply_name":    admin["name"],     //添加人名
		"review_at":     0,                 //修改时间
		"review_uid":    0,                 //修改人uid
		"review_name":   "",                //修改人名
		"remark":        remark,            //备注
		"review_remark": "",                //审核备注
	}
	query, _, _ := dialect.Insert("tbl_agency_transfer_apply").Rows(record).ToSQL()
	fmt.Println(query)
	_, err := meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	return nil
}

//MemberTransferReview 团队转代申请审核
func MemberTransferReview(id, reviewRemark string, status int, admin map[string]string) error {

	ex := g.Ex{
		"id": id,
	}
	var st int
	query, _, _ := dialect.From("tbl_agency_transfer_apply").Select("status").Where(ex).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Get(&st, query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	if st == status {
		return errors.New(helper.NoDataUpdate)
	}

	record := g.Record{
		"status":        status,
		"review_at":     time.Now().Unix(),
		"review_uid":    admin["id"],
		"review_name":   admin["name"],
		"review_remark": reviewRemark,
	}
	query, _, _ = dialect.Update("tbl_agency_transfer_apply").Set(record).Where(ex).ToSQL()
	fmt.Println(query)
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	return nil
}

//MemberTransferDelete 团队转代申请删除
func MemberTransferDelete(id string, admin map[string]string) error {

	ex := g.Ex{
		"id": id,
	}
	var st int
	query, _, _ := dialect.From("tbl_agency_transfer_apply").Select("status").Where(ex).ToSQL()
	fmt.Println(query)
	err := meta.MerchantDB.Get(&st, query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	if st == 4 {
		return errors.New(helper.NoDataUpdate)
	}

	record := g.Record{
		"status": 4,
	}
	query, _, _ = dialect.Update("tbl_agency_transfer_apply").Set(record).Where(ex).ToSQL()
	fmt.Println(query)
	_, err = meta.MerchantDB.Exec(query)
	if err != nil {
		return pushLog(err, helper.DBErr)
	}

	return nil
}
