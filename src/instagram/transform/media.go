package transform

import "github.com/feelthecode/instagramrobot/src/instagram/response"

type Owner struct {
	Id                string `json:"id"`
	ProfilePictureURL string `json:"profile_pic_url"`
	Username          string `json:"username"`
	Followers         uint64 `json:"followers_count"`
	IsPrivate         bool   `json:"is_private"`
	IsVerified        bool   `json:"is_verified"`
}

type MediaItem struct {
	Id        string `json:"id"`
	Shortcode string `json:"shortcode"`
	Type      string `json:"type"`
	IsVideo   bool   `json:"is_video"`
	Url       string `json:"url"`
}

type Media struct {
	Id        string      `json:"id"`
	Shortcode string      `json:"shortcode"`
	Type      string      `json:"type"`
	Comments  uint64      `json:"comments_count"`
	Likes     uint64      `json:"likes_count"`
	Caption   string      `json:"caption"`
	IsVideo   bool        `json:"is_video"`
	Url       string      `json:"url"`
	Items     []MediaItem `json:"items"`
	TakenAt   int64       `json:"taken_at"` // Timestamp
}

func FromEmbedResponse(embed response.EmbedResponse) Media {
	media := Media{
		Id:        embed.Media.Id,
		Shortcode: embed.Media.Shortcode,
		Type:      embed.Media.Type,
		Comments:  embed.Media.Comments.Count,
		Likes:     embed.Media.Likes.Count,
		Url:       embed.ExtractMediaURL(),
		TakenAt:   embed.Media.TakenAt.Unix(),
		IsVideo:   embed.IsVideo(),
		Caption:   embed.GetCaption(),
	}

	for _, item := range embed.Media.SliderItems.Edges {
		media.Items = append(media.Items, MediaItem{
			Id:        item.Node.Id,
			Shortcode: item.Node.Shortcode,
			Type:      item.Node.Type,
			Url:       item.Node.ExtractMediaURL(),
		})
	}

	return media
}
