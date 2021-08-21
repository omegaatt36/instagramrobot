package response

import "github.com/feelthecode/instagramrobot/src/types"

type Dimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type Resource struct {
	Width  int    `json:"config_width"`
	Height int    `json:"config_height"`
	Src    string `json:"src"`
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

type OwnerTimeline struct {
	Count uint64 `json:"count"`
	Edges []struct {
		Node struct {
			Id                 string     `json:"string"`
			ThumbnailSrc       string     `json:"thumbnail_src"`
			ThumbnailResources []Resource `json:"thumbnail_resources"`
		} `json:"node"`
	} `json:"edges"`
}

type Owner struct {
	Id                string        `json:"id"`
	ProfilePictureURL string        `json:"profile_pic_url"`
	Username          string        `json:"username"`
	HasPublicStory    bool          `json:"has_public_story"`
	IsPrivate         bool          `json:"is_private"`
	IsVerified        bool          `json:"is_verified"`
	FollowedBy        WithCount     `json:"edge_followed_by"`
	Timeline          OwnerTimeline `json:"edge_owner_to_timeline_media"`
}

type SliderItems struct {
	Edges []struct {
		Node struct {
			Id               string     `json:"id"`
			Shortcode        string     `json:"shortcode"`
			Type             string     `json:"__typename"`
			ProductType      string     `json:"product_type"`
			Dimensions       Dimensions `json:"dimensions"`
			DisplayURL       string     `json:"display_url"`
			DisplayResources []Resource `json:"display_resources"`

			IsVideo        bool   `json:"is_video"`
			Title          string `json:"title"`
			VideoURL       string `json:"video_url"`
			VideoViewCount uint64 `json:"video_view_count"`

			// clips_music_attribution_info
			// media_overlay_info
			// sharing_friction_info
		} `json:"node"`
	} `json:"edges"`
}

type ShortcodeMedia struct {
	Type             string     `json:"__typename"`
	Id               string     `json:"id"`
	Shortcode        string     `json:"shortcode"`
	ProductType      string     `json:"product_type"`
	TakenAt          types.Time `json:"taken_at_timestamp"`
	CommenterCount   uint64     `json:"commenter_count"`
	Comments         WithCount  `json:"edge_media_to_comment"`
	Likes            WithCount  `json:"edge_liked_by"`
	Dimensions       Dimensions `json:"dimensions"`
	DisplayURL       string     `json:"display_url"`
	DisplayResources []Resource `json:"display_resources"`

	Caption Caption `json:"edge_media_to_caption"`
	Owner   Owner   `json:"owner"`

	IsVideo        bool   `json:"is_video"`
	Title          string `json:"title"`
	VideoURL       string `json:"video_url"`
	VideoViewCount uint64 `json:"video_view_count"`

	SliderItems SliderItems `json:"edge_sidecar_to_children"`

	// clips_music_attribution_info
	// media_overlay_info
	// sharing_friction_info

	// edge_media_to_sponsor_user
	// is_affiliate
	// is_paid_partnership
	// location
	// coauthor_producers
}

type EmbedResponse struct {
	Media ShortcodeMedia `json:"shortcode_media"`
}

func (s EmbedResponse) IsEmpty() bool {
	return s.Media.Id == ""
}

func (s EmbedResponse) IsVideo() bool {
	return s.Media.IsVideo
}

func (s EmbedResponse) GetCaption() string {
	return s.Media.Caption.Edges[0].Node.Text
}

func (s EmbedResponse) ExtractMediaURL() string {
	if s.Media.IsVideo {
		return s.Media.VideoURL
	}
	return s.Media.DisplayURL
}
