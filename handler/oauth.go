package handler

import (
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"

	//"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iam/v1"

	//"cloud.google.com/go/iam/admin/apiv1"
	//adminpb "google.golang.org/genproto/googleapis/iam/admin/v1"
	//"google.golang.org/api/iterator"

	//"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	conf = oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"openid", "email", "profile", iam.CloudPlatformScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://www.googleapis.com/oauth2/v4/token",
		},
	}
)

const OAuthCookieName = "OAuthState"

//http.SetCookie(w,sc)
func redirectLogin(w http.ResponseWriter, r *http.Request, uri string) {

	state := uuid.NewV4().String()
	sc := http.Cookie{
		Name:   OAuthCookieName,
		Value:  state,
		MaxAge: 60,
		Path:   "/",
	}
	http.SetCookie(w, &sc)

	conf.RedirectURL = fmt.Sprintf("https://%s%s", "chic-shizuoka.appspot.com", uri)
	url := conf.AuthCodeURL(state)

	http.Redirect(w, r, url, 302)
}

func setToken(w http.ResponseWriter, r *http.Request) error {

	sc, err := r.Cookie(OAuthCookieName)
	if err != nil {
		return err
	}

	if sc.Value != r.FormValue("state") {
		return fmt.Errorf("Cookie State Value Error")
	}

	code := r.FormValue("code")
	ctx := r.Context()
	tok, err := conf.Exchange(ctx, code)
	if !tok.Valid() {
		return fmt.Errorf("Token Valid Error")
	}
	/*
		client := conf.Client(ctx, tok)
		cloudresourcemanagerService, err := cloudresourcemanager.New(client)
		if err != nil {
			log.Println("New:", err)
			return err
		}
	*/

	service, _ := iam.New(conf.Client(ctx, tok))
	resource := "projects/*/serviceAccounts/*"
	resp, err := service.Projects.ServiceAccounts.GetIamPolicy(resource).Do()
	if err != nil {
		log.Println("GetIamPolicy:", err)
		return err
	}

	for _, bind := range resp.Bindings {
		log.Println("Role: " + bind.Role)
		log.Println("Member: ")
		for _, mem := range bind.Members {
			log.Println(mem)
		}
	}

	/*
		rb := &cloudresourcemanager.GetIamPolicyRequest{}
		resp, err := cloudresourcemanagerService.Projects.GetIamPolicy(resource, rb).Context(ctx).Do()
		if err != nil {
			log.Println("GetIamPolicy:", err)
			return err
		}

	*/

	/*
		log.Println("IamClient")
		c, err := admin.NewIamClient(ctx)
		if err != nil {
			return err
		}

		req := adminpb.ListRolesRequest{
			Parent: "projects/chic-shizuoka",
		}

		log.Println("ListRoles")
		it := c.ListRoles(ctx, &req)
		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			fmt.Println(resp)
		}
	*/
	return nil
}
