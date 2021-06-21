package go_pocket_sdk

const (
	Add         = "add"
	Archive     = "archive"
	ReAdd       = "readd"
	Favorite    = "favorite"
	UnFavorite  = "unfavorite"
	Delete      = "delete"
	TagsAdd     = "tags_add"
	TagsRemove  = "tags_remove"
	TagsReplace = "tags_replace"
	TagsClear   = "tags_clear"
	TagRename   = "tag_rename"
	TagDelete   = "tag_delete"
)

type Action struct {
	Name   string `json:"action"`
	ItemId string `json:"item_id"`
	RefId  string `json:"ref_id,omitempty"`
	Tags   string `json:"tags,omitempty"`
	Time   int64  `json:"time,omitempty"`
	Title  string `json:"title,omitempty"`
	Url    string `json:"url,omitempty"`
	OldTag string `json:"old_tag,omitempty"`
	NewTag string `json:"new_tag,omitempty"`
}
