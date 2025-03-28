package instagram

import (
	"strconv"
	"time"
)

// Dimensions represents the height and width of a media item.
type Dimensions struct {
	// Height in pixels.
	Height int `json:"height"`
	// Width in pixels.
	Width int `json:"width"`
}

// Resource represents a specific resolution variant of a media item.
type Resource struct {
	// Width of the resource in pixels.
	Width int `json:"config_width"`
	// Height of the resource in pixels.
	Height int `json:"config_height"`
	// Src is the direct URL to this specific resource variant.
	Src string `json:"src"`
}

// Caption contains the text caption associated with an Instagram post.
type Caption struct {
	// Edges contain the nodes holding the caption text. Usually only one edge/node.
	Edges []struct {
		// Node contains the actual caption text.
		Node struct {
			// Text is the raw caption content.
			Text string `json:"text"`
		} `json:"node"`
	} `json:"edges"`
}

// WithCount represents a field that primarily holds a count.
type WithCount struct {
	// Count is the numerical value.
	Count uint64 `json:"count"`
}

// OwnerTimeline represents a summary of the owner's media timeline.
type OwnerTimeline struct {
	// Count is the total number of public posts by the owner.
	Count uint64 `json:"count"`
	// Edges contain nodes representing recent media items on the owner's timeline.
	Edges []struct {
		// Node contains details of a media item thumbnail.
		Node struct {
			Id string `json:"string"` // Note: JSON key is literally "string"
			// ThumbnailSrc is the URL of a thumbnail image.
			ThumbnailSrc string `json:"thumbnail_src"`
			// ThumbnailResources provides different resolution thumbnails.
			ThumbnailResources []Resource `json:"thumbnail_resources"`
		} `json:"node"`
	} `json:"edges"`
}

// Owner represents the Instagram user who created the post.
type Owner struct {
	// Id is the unique numerical ID of the user.
	Id string `json:"id"`
	// ProfilePictureURL is the URL to the user's profile picture.
	ProfilePictureURL string `json:"profile_pic_url"`
	// Username is the user's unique username.
	Username string `json:"username"`
	// HasPublicStory indicates if the user has a currently viewable public story.
	HasPublicStory bool `json:"has_public_story"`
	// IsPrivate indicates if the user's account is private.
	IsPrivate bool `json:"is_private"`
	// IsVerified indicates if the user account has a verified badge.
	IsVerified bool `json:"is_verified"`
	// Followers contains the count of users following this owner.
	Followers WithCount `json:"edge_followed_by"`
	// Timeline provides a glimpse into the owner's recent posts.
	Timeline OwnerTimeline `json:"edge_owner_to_timeline_media"`
}

// SliderItemNode represents a single media item within a multi-item post (carousel).
type SliderItemNode struct {
	// Id is the unique ID of this specific media item.
	Id string `json:"id"`
	// ShortCode is the shortcode of the parent post.
	ShortCode string `json:"shortcode"`
	// Type indicates the type of media (e.g., "GraphImage", "GraphVideo").
	Type string `json:"__typename"`
	// ProductType likely relates to Instagram's product tagging features.
	ProductType string `json:"product_type"`
	// Dimensions are the pixel dimensions of this media item.
	Dimensions Dimensions `json:"dimensions"`
	// DisplayURL is a URL to a version of the media (resolution might vary).
	DisplayURL string `json:"display_url"`
	// DisplayResources lists available resolutions/versions of the media.
	DisplayResources []Resource `json:"display_resources"`
	// IsVideo is true if this item is a video.
	IsVideo bool `json:"is_video"`
	// Title might contain a title, especially for videos (often empty for images).
	Title string `json:"title"`
	// VideoURL is the direct URL to the video file (if IsVideo is true).
	VideoURL string `json:"video_url"`
	// VideoViewCount is the number of views for this video item.
	VideoViewCount uint64 `json:"video_view_count"`
}

// ExtractMediaURL returns the most appropriate media URL (VideoURL if video, DisplayURL otherwise).
func (s SliderItemNode) ExtractMediaURL() string {
	if s.IsVideo {
		return s.VideoURL
	}
	return s.DisplayURL
}

// SliderItems contains the list of individual media items in a carousel post.
type SliderItems struct {
	// Edges contain the nodes, each representing a SliderItemNode.
	Edges []struct {
		// Node holds the actual data for a single media item in the carousel.
		Node SliderItemNode `json:"node"`
	} `json:"edges"`
}

// Media represents the core data structure for an Instagram post extracted from embedded data.
type Media struct {
	// ID is the unique ID of the post.
	ID string `json:"id"`
	// ShortCode is the unique short identifier of the post (e.g., "CqXv_...).
	ShortCode string `json:"shortcode"`
	// Type indicates the overall type of the post (e.g., "GraphImage", "GraphVideo", "GraphSidecar").
	Type string `json:"__typename"`
	// ProductType likely relates to Instagram's product tagging features.
	ProductType string `json:"product_type"`
	// TakenAt is the timestamp when the media was published.
	TakenAt Time `json:"taken_at_timestamp"`
	// CommenterCount seems deprecated or unused in recent GQL data.
	CommenterCount uint64 `json:"commenter_count"` // Might be deprecated
	// Comments contains the total count of comments on the post.
	Comments WithCount `json:"edge_media_to_comment"`
	// Likes contains the total count of likes on the post.
	Likes WithCount `json:"edge_liked_by"`
	// Dimensions are the pixel dimensions of the primary media item (if not a sidecar).
	Dimensions Dimensions `json:"dimensions"`
	// DisplayURL is a URL to a version of the primary media (resolution might vary).
	DisplayURL string `json:"display_url"`
	// DisplayResources lists available resolutions/versions of the primary media.
	DisplayResources []Resource `json:"display_resources"`
	// Caption holds the text caption of the post.
	Caption Caption `json:"edge_media_to_caption"`
	// Owner contains information about the user who published the post.
	Owner Owner `json:"owner"`
	// IsVideo is true if the primary media item is a video (not applicable for sidecar).
	IsVideo bool `json:"is_video"`
	// Title might contain a title, especially for videos.
	Title string `json:"title"`
	// VideoURL is the direct URL to the video file (if IsVideo is true and not a sidecar).
	VideoURL string `json:"video_url"`
	// VideoViewCount is the number of views for the video (if IsVideo is true).
	VideoViewCount uint64 `json:"video_view_count"`
	// SliderItems contains the individual items if the post is a carousel (Type == "GraphSidecar").
	SliderItems SliderItems `json:"edge_sidecar_to_children"`
}

// EmbedResponse is the root structure holding the Media data parsed from the GQL JSON.
type EmbedResponse struct {
	// Media contains the main post data.
	Media Media `json:"shortcode_media"`
}

// IsEmpty checks if the embedded Media data was successfully parsed (based on ID).
func (s EmbedResponse) IsEmpty() bool {
	return s.Media.ID == ""
}

// IsVideo checks if the primary media item in the response is a video.
// Note: For carousels (sidecar), this reflects the *first* item's type inaccurately. Check SliderItems instead.
func (s EmbedResponse) IsVideo() bool {
	return s.Media.IsVideo
}

// GetCaption extracts the caption text from the Media data.
// It falls back to the Title if no caption edges are found.
func (s EmbedResponse) GetCaption() string {
	if len(s.Media.Caption.Edges) > 0 {
		return s.Media.Caption.Edges[0].Node.Text
	}
	// Fallback for potential cases where caption is in title (less common now)
	return s.Media.Title
}

// ExtractMediaURL returns the most appropriate URL for the primary media item.
// Note: For carousels (sidecar), this returns the DisplayURL of the *post* itself, not an individual item.
// Use SliderItemNode.ExtractMediaURL for carousel items.
func (s EmbedResponse) ExtractMediaURL() string {
	if s.Media.IsVideo {
		return s.Media.VideoURL
	}
	return s.Media.DisplayURL
}

// Time defines a custom type for handling Unix timestamps encoded as numbers in JSON.
type Time time.Time

// MarshalJSON converts the Time type to its Unix timestamp integer representation for JSON encoding.
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

// UnmarshalJSON converts a Unix timestamp integer from JSON data into the Time type.
func (t *Time) UnmarshalJSON(s []byte) (err error) {
	r := string(s)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	// Dereference the pointer and assign the converted time.Time value.
	*(*time.Time)(t) = time.Unix(q, 0)
	return nil
}

// Unix returns t as a Unix time (seconds since epoch).
func (t Time) Unix() int64 {
	return time.Time(t).Unix()
}

// Time returns the underlying time.Time value in UTC.
func (t Time) Time() time.Time {
	return time.Time(t).UTC()
}

// String returns the formatted string representation of the time.Time value.
func (t Time) String() string {
	return t.Time().String()
}
