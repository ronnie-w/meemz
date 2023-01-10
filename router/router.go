package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meemz/activities"
	"github.com/meemz/authentication"
	"github.com/meemz/contentfetch"
	"github.com/meemz/parser"
	"github.com/meemz/profile"
	"github.com/meemz/regommend"
	"github.com/meemz/upload"
)

func Routes() *mux.Router {
	mux := mux.NewRouter()

	templates := map[string]func(http.ResponseWriter, *http.Request){
		"/":       parser.Parser("home"),
		"/signup": parser.Parser("signup"),
		"/login":  parser.Parser("login"),
		// "/terms":           parser.Terms,
		"/verify":  parser.Parser("verify"),
		"/profile": parser.Parser("profile"),
		// "/private_posts":   parser.PrivatePosts,
		// "/profile_edit":    parser.ProfileEdit,
		"/create":             parser.Parser("create"),
		"/search":             parser.Parser("search"),
		"/notifications":      parser.Parser("notifications"),
		"/comments/{file_id}": parser.Parser("comments"),
		"/forgot_password":    parser.Parser("forgot_password"),
		"/verify_passcode":    parser.Parser("verify_passcode"),
		"/password_reset":     parser.Parser("password_reset"),

		// "/user/{username}":         parser.Account,
		// "/public_stats/{img_name}": parser.PublicStats,
		// "/image_stats/{img_name}":  parser.ImageStats,

		"/manifest.json": parser.Manifest,
		"/sw_init.js":    parser.ServiceWorkerInit,
		"/sw.js":         parser.ServiceWorker,

		"/.well-known/assetslinks.json": parser.AssetsLink,
	}

	end_points := map[string]func(http.ResponseWriter, *http.Request){
		//regommendation router
		"/main_content":       regommend.FetchMeemz,
		"/main_veemz_content": regommend.FetchVeemz,

		//authentication router
		"/signup_auth":       authentication.Signup,
		"/wrong_mail":        authentication.DeleteUser,
		"/fetch_user":        authentication.FetchUser,
		"/verify_auth":       authentication.Verify,
		"/login_auth":        authentication.Login,
		"/pass_send_code":    authentication.SendVerificationCode,
		"/pass_verify_code":  authentication.ConfirmVerificationCode,
		"/pass_new_password": authentication.PasswordReset,

		//upload router
		"/upload_meemz":        upload.UploadMeemz,
		"/upload_veemz":        upload.UploadVeemz,
		"/update_meemz_config": upload.UpdateMeemzConfig,
		"/update_veemz_config": upload.UpdateVeemzConfig,
		"/generate_new_id":     upload.GenerateNewId,

		//profile router
		"/my_uploads":              contentfetch.MyUploads,
		"/fetch_user_uploads":      contentfetch.UsersUploads,
		"/profile_img_upload":      profile.ProfileUpload,
		"/profile_change_username": profile.PostUsernameToDb,
		"/profile_change_email":    profile.PostEmailToDb,
		"/profile_change_bio":      profile.PostBioToDb,
		"/profile_update_img":      profile.ProfileUpdateImg,
		"/subscribe":               profile.Subscribe,
		"/unsubscribe":             profile.UnSubscribe,
		"/fetch_subs":              profile.FetchSubs,
		"/delete_post":             profile.DeletePost,

		//search router
		"/search_meemz": contentfetch.SearchMeemz,
		"/search_users": contentfetch.SearchUsers,
		"/search_tags":  contentfetch.SearchTags,

		//activities router
		// "/post_reaction":   activities.PostReaction,
		// "/delete_reaction": activities.DeleteReaction,
		// "/post_report":     activities.PostReport,
		// "/delete_report":   activities.DeleteReport,

		"/fetch_comments":          activities.FetchComments,
		"/fetch_replies":           activities.FetchReplies,
		"/post_comment":            activities.PostComment,
		"/post_reply":              activities.PostReply,
		"/post_comment_reaction":   activities.PostCommentReplyLike,
		"/delete_comment":          activities.DeleteComment,
		"/delete_reply":            activities.DeleteReply,
		"/delete_comment_reaction": activities.DeleteCommentReplyLike,

		"/fetch_notifications":  activities.FetchNotifications,
		"/delete_notifications": activities.DeleteNotifications,
		"/ws":                   activities.ClientHandler,
	}

	for route, http_func := range templates {
		mux.HandleFunc(route, http_func)
	}

	for route, http_func := range end_points {
		mux.HandleFunc(route, http_func)
	}

	fileServer := http.FileServer(http.Dir("static/"))
	mux.PathPrefix("/").Handler(http.StripPrefix("/static", fileServer))

	return mux
}
