package model

func MemberClosureInsert(nodeID, targetID string) string {

	t := "SELECT ancestor, " + nodeID + ",prefix, lvl+1 FROM tbl_members_tree WHERE prefix='" + meta.Prefix + "' and descendant = " + targetID + " UNION SELECT " + nodeID + "," + nodeID + "," + "'" + meta.Prefix + "'" + ",0"
	query := "INSERT INTO tbl_members_tree (ancestor, descendant,prefix,lvl) (" + t + ")"

	return query
}
