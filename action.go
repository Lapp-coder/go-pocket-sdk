package go_pocket_sdk

const (
	ActionAdd         = "add"
	ActionArchive     = "archive"
	ActionReAdd       = "readd"
	ActionFavorite    = "favorite"
	ActionUnFavorite  = "unfavorite"
	ActionDelete      = "delete"
	ActionTagsAdd     = "tags_add"
	ActionTagsRemove  = "tags_remove"
	ActionTagsReplace = "tags_replace"
	ActionTagsClear   = "tags_clear"
	ActionTagRename   = "tag_rename"
	ActionTagDelete   = "tag_delete"
)

type Action struct {
	Name   string `json:"action"`
	ItemID string `json:"item_id"`
	RefId  string `json:"ref_id,omitempty"`
	Tags   string `json:"tags,omitempty"`
	Time   int64  `json:"time,omitempty"`
	Title  string `json:"title,omitempty"`
	Url    string `json:"url,omitempty"`
	OldTag string `json:"old_tag,omitempty"`
	NewTag string `json:"new_tag,omitempty"`
}
