package model

import (
	"fmt"
	"github.com/doug-martin/goqu/v9/exp"
	"merchant/contrib/helper"
	"strconv"
	"time"

	g "github.com/doug-martin/goqu/v9"
	"github.com/shopspring/decimal"
)

// 是否允许修改下级返水比例
func MemberRebateEnableMod(enable bool) error {

	key := fmt.Sprintf("%s:rebate:enablemod", meta.Prefix)
	// 允许修改下级返水比例
	if enable {
		pipe := meta.MerchantRedis.TxPipeline()
		defer pipe.Close()

		pipe.Set(ctx, key, "1", 100*time.Hour)
		pipe.Persist(ctx, key)

		_, err := pipe.Exec(ctx)
		if err != nil {
			return pushLog(err, helper.RedisErr)
		}
	} else { //禁止修改下级返水比例
		cmd := meta.MerchantRedis.Del(ctx, key)
		fmt.Println(cmd.String())
		err := cmd.Err()
		if err != nil {
			return pushLog(err, helper.RedisErr)
		}
	}

	return nil
}

func MemberRebateScale() MemberRebateResult_t {
	return meta.VenueRebate
}

func MemberRebateCmp(lower, own MemberRebateResult_t) bool {

	if own.QP.LessThan(lower.QP) {
		return false
	}
	if own.ZR.LessThan(lower.ZR) {
		return false
	}
	if own.TY.LessThan(lower.TY) {
		return false
	}
	if own.DJ.LessThan(lower.DJ) {
		return false
	}
	if own.DZ.LessThan(lower.DZ) {
		return false
	}
	if own.CP.LessThan(lower.CP) {
		return false
	}
	if own.FC.LessThan(lower.FC) {
		return false
	}
	if own.BY.LessThan(lower.BY) {
		return false
	}
	if own.CGHighRebate.LessThan(lower.CGHighRebate) {
		return false
	}
	if own.CGOfficialRebate.LessThan(lower.CGOfficialRebate) {
		return false
	}

	return true
}

func MemberRebateUpdateCache1(uid string, mr MemberRebateResult_t) error {

	key := fmt.Sprintf("%s:m:rebate:%s", meta.Prefix, uid)
	vals := []interface{}{"zr", mr.ZR.StringFixed(1), "qp", mr.QP.StringFixed(1), "ty", mr.TY.StringFixed(1), "dj", mr.DJ.StringFixed(1), "dz", mr.DZ.StringFixed(1), "cp", mr.CP.StringFixed(1), "fc", mr.FC.StringFixed(1), "by", mr.BY.StringFixed(1), "cg_high_rebate", mr.CGHighRebate.StringFixed(2), "cg_official_rebate", mr.CGOfficialRebate.StringFixed(2)}

	pipe := meta.MerchantRedis.Pipeline()
	pipe.Del(ctx, key)
	pipe.HMSet(ctx, key, vals...)
	pipe.Persist(ctx, key)
	_, err := pipe.Exec(ctx)
	pipe.Close()

	return err
}

func MemberRebateUpdateCache2(uid string, mr MemberRebate) error {

	key := fmt.Sprintf("%s:m:rebate:%s", meta.Prefix, uid)
	vals := []interface{}{"zr", mr.ZR, "qp", mr.QP, "ty", mr.TY, "dj", mr.DJ, "dz", mr.DZ, "cp", mr.CP, "fc", mr.FC, "by", mr.BY, "cg_high_rebate", mr.CgHighRebate, "cg_official_rebate", mr.CgOfficialRebate}

	pipe := meta.MerchantRedis.Pipeline()
	pipe.Del(ctx, key)
	pipe.HMSet(ctx, key, vals...)
	pipe.Persist(ctx, key)
	_, err := pipe.Exec(ctx)
	pipe.Close()

	return err
}

func MemberRebateFindOne(uid string) (MemberRebateResult_t, error) {

	data := MemberRebate{}
	res := MemberRebateResult_t{}

	t := dialect.From("tbl_member_rebate_info")
	query, _, _ := t.Select(colsMemberRebate...).Where(g.Ex{"uid": uid}).Limit(1).ToSQL()
	err := meta.MerchantDB.Get(&data, query)
	if err != nil {
		return res, pushLog(err, helper.DBErr)
	}

	res.ZR, _ = decimal.NewFromString(data.ZR)
	res.QP, _ = decimal.NewFromString(data.QP)
	res.TY, _ = decimal.NewFromString(data.TY)
	res.DJ, _ = decimal.NewFromString(data.DJ)
	res.DZ, _ = decimal.NewFromString(data.DZ)
	res.CP, _ = decimal.NewFromString(data.CP)
	res.FC, _ = decimal.NewFromString(data.FC)
	res.BY, _ = decimal.NewFromString(data.BY)
	res.CGOfficialRebate, _ = decimal.NewFromString(data.CgOfficialRebate)
	res.CGHighRebate, _ = decimal.NewFromString(data.CgHighRebate)

	res.ZR = res.ZR.Truncate(1)
	res.QP = res.QP.Truncate(1)
	res.TY = res.TY.Truncate(1)
	res.DJ = res.DJ.Truncate(1)
	res.DZ = res.DZ.Truncate(1)
	res.CP = res.CP.Truncate(1)
	res.FC = res.FC.Truncate(1)
	res.BY = res.BY.Truncate(1)

	res.CGOfficialRebate = res.CGOfficialRebate.Truncate(2)
	res.CGHighRebate = res.CGHighRebate.Truncate(2)

	return res, nil
}

// RebatePersonalReport 个人返水记录
func RebatePersonalReport(username string, startTime, endTime string, page, pageSize int) (PersonalRebateReportData, error) {

	data := PersonalRebateReportData{}
	dayMap := map[string]MemberPersonalRebate{}

	startAt, err := helper.TimeToLoc(startTime, loc)
	if err != nil {
		return data, pushLog(err, helper.DateTimeErr)
	}

	endAt, err := helper.TimeToLoc(endTime, loc)
	if err != nil {
		return data, pushLog(err, helper.DateTimeErr)
	}

	totalDays := RebateCountDays(startAt, endAt)
	//处理开始时间和结束时间，根据翻页参数把开始时间和结束时间计算出来，只获取当前页的数据即可
	fmt.Printf("before count startAt:%d,endAt:%d,page:%d,%d", startAt, endAt, page, pageSize)
	err = RebateDealEndTime(&startAt, &endAt, page, pageSize)
	if err != nil {
		fmt.Println("根据开始时间,结束时间,页数，页容量获取开始时间和结束时间出错")
		return data, err
	}

	fmt.Printf("after count startAt:%d,endAt:%d", startAt, endAt)

	//根据开始时间，和结束时间，pageSize返回一个slice
	sliceEndTime, err := RebateGetResSliceEndTime(startAt, endAt, pageSize)
	if err != nil {
		fmt.Println("根据开始时间和结束时间获取slice出错")
		return data, err
	}

	//制造一个map封装所有的返回参数
	var memberRebate []RebateReportItem
	ex := g.Ex{
		"prefix":      meta.Prefix,
		"username":    username,
		"report_time": g.Op{"between": exp.NewRangeVal(startAt, endAt)},
	}
	query, _, _ := dialect.From("tbl_report_member_rebate").Select(colRebateReport...).Where(ex).ToSQL()
	err = meta.ReportDB.Select(&memberRebate, query)
	fmt.Println(query)
	if err != nil {
		return data, pushLog(fmt.Errorf("%s,[%s]", err.Error(), query), helper.DBErr)
	}

	dayMap = GetInitResMap(sliceEndTime)
	for _, v := range memberRebate {

		//string转int64的时间戳
		reportTime, err := strconv.ParseInt(v.ReportTime, 10, 64)
		if err != nil {
			return data, pushLog(err, helper.DateTimeErr)
		}
		tm := time.Unix(reportTime, 0)
		day := tm.Format("2006-01-02")
		rebate := dayMap[day]
		amount, _ := decimal.NewFromString(v.RebateAmount)
		switch v.CashType {
		case helper.TransactionRebateCasino:
			rebate.ZR = rebate.ZR.Add(amount)
		case helper.TransactionRebateLottery:
			rebate.CP = rebate.CP.Add(amount)
		case helper.TransactionRebateSport:
			rebate.TY = rebate.TY.Add(amount)
		case helper.TransactionRebateDesk:
			rebate.QP = rebate.QP.Add(amount)
		case helper.TransactionRebateESport:
			rebate.DJ = rebate.DJ.Add(amount)
		case helper.TransactionRebateCockFighting:
			rebate.FC = rebate.FC.Add(amount)
		case helper.TransactionRebateFishing:
			rebate.BY = rebate.BY.Add(amount)
		case helper.TransactionRebateLott:
			rebate.DZ = rebate.DZ.Add(amount)
			//case helper.TransactionRebate:
			//	rebate.TotalRebate = rebate.TotalRebate.Add(validBet)
		}
		dayMap[day] = rebate
	}
	var resData []MemberPersonalRebate

	for _, v := range sliceEndTime {
		rebate := dayMap[v]
		//计算汇总
		rebate.TotalRebate = rebate.ZR.Add(rebate.CP).Add(rebate.TY).Add(rebate.QP).Add(rebate.DJ).Add(rebate.FC).
			Add(rebate.BY).Add(rebate.DZ)
		rebate.RebateDate = v
		resData = append(resData, rebate)
	}
	data.D = resData
	data.T = totalDays
	data.S = uint16(pageSize)
	return data, nil
}

//RebateDealEndTime 据开始时间，结束时间，当前页数，翻页容量。修改开始时间和结束时间
func RebateDealEndTime(startAt, endAt *int64, page, pageSize int) error {

	//如果是第一页就只计算当前页的开始时间和解释时间
	if page == 1 {
		//offset := (page - 1) * pageSize
		//计算结束时间
		endAtTime := time.Unix(*endAt, 0)
		startAtTimeBySize := endAtTime.AddDate(0, 0, -pageSize).Unix()
		if *startAt < startAtTimeBySize {
			tm := time.Unix(startAtTimeBySize, 0)
			//获取当天的00:00:00
			format := tm.Format("2006-01-02") + " 00:00:00"
			startAt64, err := helper.TimeToLoc(format, loc)
			if err != nil {
				return pushLog(err, helper.DateTimeErr)
			}

			*startAt = startAt64
		}
	} else {
		endAtTime := time.Unix(*endAt, 0)
		//计算翻页后的开始时间和结束时间
		offset := (page - 1) * pageSize
		endAtTimeTimeBySize := endAtTime.AddDate(0, 0, -offset)
		//如果计算的结束时间在开始时间以前，则报错
		if endAtTimeTimeBySize.Unix() < *startAt {
			fmt.Print("根据分页计算开始时间和结束时间错误", startAt, endAt, page, pageSize)
			return pushLog(fmt.Errorf("%d,[%d]", endAtTimeTimeBySize, startAt), helper.ParamErr)
		}

		//如果计算的结束时间在结束时间以后，则报错
		if *endAt < endAtTimeTimeBySize.Unix() {
			fmt.Print("根据分页计算开始时间和结束时间错误", startAt, endAt, page, pageSize)
			return pushLog(fmt.Errorf("endAt:%d,endAtTimeTimeBySize:[%d]", endAt, endAtTimeTimeBySize), helper.ParamErr)

		}

		format := endAtTimeTimeBySize.Format("2006-01-02") + " 23:59:59"
		endAtString, err := helper.TimeToLoc(format, loc)
		if err != nil {
			return pushLog(err, helper.DateTimeErr)
		}
		*endAt = endAtString
		endTimeEnd := time.Unix(endAtString, 0)
		startAtTimeBySize := endTimeEnd.AddDate(0, 0, -pageSize).Unix()
		if *startAt < startAtTimeBySize {
			tm := time.Unix(startAtTimeBySize, 0)
			//获取当天的00:00:00
			format := tm.Format("2006-01-02") + " 00:00:00"
			startAt64, err := helper.TimeToLoc(format, loc)
			if err != nil {
				return pushLog(err, helper.DateTimeErr)
			}

			*startAt = startAt64
		}

	}
	if *startAt >= *endAt {
		return pushLog(fmt.Errorf("startAt:%d,endAt:[%d]", startAt, endAt), helper.QueryTimeRangeErr)
	}

	return nil
}

// RebateGetResSliceEndTime 根据开始时间，结束时间，当前页数，翻页容量。修改开始时间和结束时间
func RebateGetResSliceEndTime(startAt, endAt int64, pageSize int) ([]string, error) {

	//如果是第一页就只计算当前页的开始时间和解释时间
	var daySlice []string
	tm := time.Unix(startAt, 0)
	//获取当天的00:00:00
	format := tm.Format("2006-01-02") + " 00:00:00"
	startAt64, err := helper.TimeToLoc(format, loc)
	if err != nil {
		return daySlice, pushLog(err, helper.DateTimeErr)
	}

	for i := 0; i < pageSize; i++ {
		endAtTime := time.Unix(endAt, 0)
		endAtTimeTimeBySize := endAtTime.AddDate(0, 0, -i).Unix()
		//如果时间在开始时间以前，则停止
		if endAtTimeTimeBySize < startAt64 {
			break
		}

		tm := time.Unix(endAtTimeTimeBySize, 0)
		format := tm.Format("2006-01-02")
		daySlice = append(daySlice, format)
	}

	fmt.Printf("daySlice:%s,len:%d", daySlice, len(daySlice))
	return daySlice, nil
}

//GetInitResMap 根据开始时间，结束时间，当前页数，翻页容量。修改开始时间和结束时间
func GetInitResMap(sliceDay []string) map[string]MemberPersonalRebate {

	dayMap := make(map[string]MemberPersonalRebate) //创建集合
	for _, s := range sliceDay {
		personalRebate := MemberPersonalRebate{}
		dayMap[s] = personalRebate
	}

	fmt.Println("GetInitResMap", dayMap)
	return dayMap
}

//RebateCountDays 计算两个时间差距几天
func RebateCountDays(start, end int64) int64 {

	startTime := time.Unix(start, 0)
	endTime := time.Unix(end, 0)
	sub := int(endTime.Sub(startTime).Hours())
	days := (int64)(sub / 24)
	if (sub % 24) > 0 {
		days = days + 1
	}

	return days
}
