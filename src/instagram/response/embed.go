package response

import "github.com/feelthecode/instagramrobot/src/types"

type Dimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type DisplayResource struct {
	ConfigWidth  int    `json:"config_width"`
	ConfigHeight int    `json:"config_height"`
	Src          string `json:"src"`
}

type Caption struct {
	Edges []struct {
		Node struct {
			Text string `json:"text"`
		} `json:"node"`
	} `json:"edges"`
}

type WithCount struct {
	Count uint64 `json:"count"`
}

type Owner struct {
	Id                string    `json:"id"`
	ProfilePictureURL string    `json:"profile_pic_url"`
	Username          string    `json:"username"`
	HasPublicStory    bool      `json:"has_public_story"`
	IsPrivate         bool      `json:"is_private"`
	IsVerified        bool      `json:"is_verified"`
	FollowedBy        WithCount `json:"edge_followed_by"`
	// TODO: add more fields
}

type EmbedResponse struct {
	ShortcodeMedia struct {
		Type             string            `json:"__typename"`
		Id               string            `json:"id"`
		Shortcode        string            `json:"shortcode"`
		TakenAt          types.Time        `json:"taken_at_timestamp"`
		CommenterCount   uint64            `json:"commenter_count"`
		Comments         WithCount         `json:"edge_media_to_comment"`
		Likes            WithCount         `json:"edge_liked_by"`
		Dimensions       Dimensions        `json:"dimensions"`
		DisplayURL       string            `json:"display_url"`
		DisplayResources []DisplayResource `json:"display_resources"`

		Caption Caption `json:"edge_media_to_caption"`
		Owner   Owner   `json:"owner"`

		IsVideo        bool   `json:"is_video"`
		Title          string `json:"title"`
		ProductType    string `json:"product_type"`
		VideoURL       string `json:"video_url"`
		VideoViewCount uint64 `json:"video_view_count"`

		// clips_music_attribution_info
		// media_overlay_info
		// sharing_friction_info
		// edge_media_to_sponsor_user
		// is_affiliate
		// is_paid_partnership
		// location
		// coauthor_producers

	} `json:"shortcode_media"`
}
