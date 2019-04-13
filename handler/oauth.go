package handler

import (
	"fmt"
	"net/http"
	"os"

	"app/datastore"

	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
	"google.golang.org/api/cloudresourcemanager/v1"
	profile "google.golang.org/api/oauth2/v2"
)

var (
	conf = oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{profile.UserinfoEmailScope, cloudresourcemanager.CloudPlatformScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://www.googleapis.com/oauth2/v4/token",
		},
	}
)

const OAuthCookieName = "OAuthState"

//http.SetCookie(w,sc)
func redirectLogin(w http.ResponseWriter, r *http.Request, uri string) error {

	fx, err := datastore.GetOAuthValue(r.Context())
	if err != nil {
		return xerrors.Errorf("error GetOAuthValue(): %w", err)
	}

	conf.ClientID = fx.ClientID
	conf.ClientSecret = fx.ClientSecret

	state := uuid.NewV4().String()
	sc := http.Cookie{
		Name:   OAuthCookieName,
		Value:  state,
		MaxAge: 120,
		Path:   "/",
	}
	http.SetCookie(w, &sc)

	host := r.URL.Host
	host = "chic-shizuoka.appspot.com"
	conf.RedirectURL = fmt.Sprintf("https://%s%s", host, uri)
	url := conf.AuthCodeURL(state)

	http.Redirect(w, r, url, 302)
	return nil
}

func authorization(w http.ResponseWriter, r *http.Request, roles ...string) error {

	sc, err := r.Cookie(OAuthCookieName)
	if err != nil {
		return xerrors.Errorf("error Cookie(): %w", err)
	}

	if sc.Value != r.FormValue("state") {
		return xerrors.Errorf("error Cookie State value.")
	}

	code := r.FormValue("code")
	ctx := r.Context()
	tok, err := conf.Exchange(ctx, code)
	if !tok.Valid() {
		return xerrors.Errorf("error Token Valid.")
	}

	if roles == nil {
		//log.Println("roles target nil equals google acount.")
		return nil
	}

	client := conf.Client(ctx, tok)
	accountService, err := profile.New(client)
	if err != nil {
		return xerrors.Errorf("error New(OAuth API).: %w", err)
	}

	info, err := accountService.Tokeninfo().Do()
	if err != nil {
		return xerrors.Errorf("error Tokeninfo(OAuth API).: %w", err)
	}

	mail := info.Email
	cloudresourcemanagerService, err := cloudresourcemanager.New(client)
	if err != nil {
		return xerrors.Errorf("error New(CloudResourceManager API).: %w", err)
	}

	resource := os.Getenv("DATASTORE_PROJECT_ID")
	rb := &cloudresourcemanager.GetIamPolicyRequest{}
	resp, err := cloudresourcemanagerService.Projects.GetIamPolicy(resource, rb).Context(ctx).Do()
	if err != nil {
		return xerrors.Errorf("error GetIamPolicy(CloudResourceManager API).: %w", err)
	}

	for _, bind := range resp.Bindings {
		for _, role := range roles {
			if bind.Role == role {
				for _, mem := range bind.Members {
					if mem == "user:"+mail {
						return nil
					}
				}
			}
		}
	}

	return xerrors.Errorf("error not found roles")
}
