package config

const (

	PathDBs = "db"
	PathDBsBackup = "db/backup"

	/*------------------------------------------------------*/

	TabUsers                    = "src_users"
	TabCategories               = "src_categories"
	TabLikes                    = "src_likes"
	TabPosts                    = "posts"
	TabFiles                    = "files"
	TabTokens                   = "tokens"
	TabPostsLikes               = "posts_likes"
	TabPostsCategories          = "posts_categories"
	TabComments                 = "comments"
	TabCommentsLikes            = "comments_likes"
	TabSessions                 = "sessions"
	QueAttach                   = "attach"
	QueDetach                   = "detach"
	QueRestore                  = "restore"
	QueInsert2                  = "insert_2"
	QueInsert3                  = "insert_3"
	QueInsert4                  = "insert_4"
	QueInsert5                  = "insert_5"
	QueInsert6                  = "insert_6"
	QueSelect                   = "select"       // all posts without any sort
	QueSelectPosts              = "select_posts" // all posts without any sort
	QueSelectUsers              = "select_users"
	QueSelectPostsBy            = "select_posts_by" // all posts sorted by WHERE
	QueSelectPostsRatedBy       = "select_posts_rated_by"
	QueSelectCommentsRatedBy    = "select_comments_rated_by"
	QueSelectPostsAndCategoryBy = "select_posts_category_by"
	QueSelectCommentAndPostBy   = "select_posts_comment_by"
	QueSelectCategories         = "select_categories"
	QueSelectUserBy             = "select_user_by"
	QueSelectCategoryBy         = "select_category_by"
	QueSelectSessionBy          = "select_session_by"
	QueSelectLikeCountBy        = "select_post_like_count_by"
	QueSelectCommentLikeCountBy = "select_comment_like_count_by"
	QueSelectCommentsBy         = "select_comments_by"
	QueSelectLikeBy             = "select_like_by"
	QueDeleteBy                 = "delete_session_by"
	QueSelectCount              = "select_count"
	QueSelectLikedOrNot         = "select_liked_or_not"

	/*------------------------------------------------------*/

	FieldID         = "id"
	FieldLike       = "like"
	FieldUUID       = "uuid"
	FieldUser       = "user"
	FieldName       = "name"
	FieldPost       = "post"
	FieldBody       = "body"
	FieldRoot       = "root"
	FieldPath       = "path"
	FieldMIME       = "mime"
	FieldLiked      = "liked"
	FieldLikes      = "likes"
	FieldLogin      = "login"
	FieldTitle      = "title"
	FieldEmail      = "email"
	FieldExpire     = "expire"
	FieldDislike    = "dislike"
	FieldCreated    = "created"
	FieldComment    = "comment"
	FieldVariety    = "variety"
	FieldCategory   = "category"
	FieldCategories = "categories"
	FieldPassword   = "password"
	FieldIdVariety  = "id_variety"

	/*------------------------------------------------------*/

	KeyID           = "id"
	KeyLogin        = "login"
	KeyEmail        = "email"
	KeyUser         = "user"
	KeyPost         = "post"
	KeyComment      = "comment"
	KeyPostRated    = "post rated"
	KeyCommentRated = "comment rated"
	KeyLike         = "like"
	KeyLink         = "link"
	KeyDislike      = "dislike"
	KeyRated        = "rated"
	KeyCategory     = "category"
	KeyRate         = "rate"
	KeyObject       = "object"
	KeyGithub       = "github"
	KeyFacebook     = "facebook"
	KeyGoogle       = "google"
)

// MakeTables returns a list of tables that need to create
func MakeTables() []string {
	return []string{
		TabUsers,
		TabCategories,
		TabLikes,
		TabPosts,
		TabFiles,
		TabPostsLikes,
		TabPostsCategories,
		TabComments,
		TabCommentsLikes,
		TabSessions,
		TabTokens,
	}
}
