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
		"/":                parser.Home,
		"/bottomNav":       parser.BottomNav,
		"/layout":          parser.Layout,
		"/signup":          parser.Signup,
		"/login":           parser.Login,
		"/terms":           parser.Terms,
		"/verify":          parser.Verify,
		"/profile":         parser.Profile,
		"/private_posts":   parser.PrivatePosts,
		"/profile_edit":    parser.ProfileEdit,
		"/create":          parser.Create,
		"/creator":         parser.Creator,
		"/veemz":			parser.Veemz,
		"/search":          parser.Search,
		"/notifications":   parser.Notifications,
		"/forgot_password": parser.ForgotPassword,

		"/user/{username}":         parser.Account,
		"/public_stats/{img_name}": parser.PublicStats,
		"/image_stats/{img_name}":  parser.ImageStats,

		"/manifest.json": parser.Manifest,
		"/sw_init.js": parser.ServiceWorkerInit,
		"/sw.js": parser.ServiceWorker,

		"/.well-known/assetslinks.json": parser.AssetsLink,
	}

	end_points := map[string]func(http.ResponseWriter, *http.Request){
		//regommendation router
		"/main_content": regommend.FetchMeemz,

		//authentication router
		"/check_user":        authentication.CheckUser,
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
		"/upload_veemz": 		upload.UploadVeemz,
		"/update_meemz_config": upload.UpdateMeemzConfig,
		"/update_veemz_config": upload.UpdateVeemzConfig,
		"/upload_convo_images": upload.UploadConvoImages,

		//profile router
		"/my_uploads":              contentfetch.MyUploads,
		"/fetch_private_posts":     contentfetch.FetchMyPrivateUploads,
		"/fetch_user_uploads":      contentfetch.UsersUploads,
		"/publicize_post":          contentfetch.PublicizePost,
		"/delete_post":             contentfetch.DeletePost,
		"/profile_img_upload":      profile.ProfileUpload,
		"/profile_change_username": profile.PostUsernameToDb,
		"/profile_change_email":    profile.PostEmailToDb,
		"/profile_change_bio":      profile.PostBioToDb,
		"/profile_update_img":      profile.ProfileUpdateImg,
		"/subscribe":               profile.Subscribe,
		"/unsubscribe":             profile.UnSubscribe,
		"/fetch_subs":              profile.FetchSubs,

		//search router
		"/search_meemz": contentfetch.SearchMeemz,
		"/search_users": contentfetch.SearchUsers,
		"/search_tags":  contentfetch.SearchTags,
		//-------------stats router
		"/fetch_stats": contentfetch.StatsHandler,

		//activities router
		"/post_reaction":        activities.PostReaction,
		"/delete_reaction":      activities.DeleteReaction,
		"/post_report":          activities.PostReport,
		"/delete_report":        activities.DeleteReport,
		"/fetch_my_comments":    activities.FetchMyComments,
		"/fetch_o_comments":     activities.FetchOtherComments,
		"/post_comment":         activities.PostComment,
		"/viewed":               activities.Viewed,
		"/notifications_go":     activities.FetchNotifications,
		"/delete_notifications": activities.DeleteNotifications,
		"/notify_invite":        activities.NotifyInvite,
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
