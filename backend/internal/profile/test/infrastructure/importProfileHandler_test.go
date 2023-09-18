package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"testing"

	application_import "github.com/unmsmfisi-socialapplication/social_app/internal/profile/application/import"
	infrastructure_import "github.com/unmsmfisi-socialapplication/social_app/internal/profile/infrastructure/import"
	infrastructure_repository "github.com/unmsmfisi-socialapplication/social_app/internal/profile/infrastructure/repository"
)

func TestImportProfileHandler_ImportProfile(t *testing.T) {
    profileRepository := infrastructure_repository.NewProfileRepository()
    importProfileUseCase := application_import.NewImportProfileUseCase(profileRepository)
    handler := infrastructure_import.NewImportProfileHandler(importProfileUseCase)

	requestBody := []byte(`{
		"@context": [
			"https://www.w3.org/ns/activitystreams",
			{
				"@language": "es"

			}
		],
		"type": "Profile",
		"actor": "https://appsocial.com/sofia/",
		"name": "Perfil de Sofia",
		"object": {
			"id": "https://appsocial.com/sofia/id1234",
			"type": "Person",
			"name": "Sofia Rodriguez",
			"preferredUsername": "SofiR",
			"summary": "Amante de la naturaleza y entusiasta de la tecnología.",
			"profileImage": "https://appsocial.com/sofia/profileImage.jpg",
			"coverImage": "https://appsocial.com/sofia/coverImage.jpg",
			"endpoints": {
				"sharedInbox": "https://appsocial.com/inbox"

			},

			"importedFrom": "https://mastodon.ejemplo.com/@sofiR"
		},
		"to": [
			"https://appsocial.com/usuarios/"
		],
		"cc": "https://appsocial.com/seguidores/sofia"
	}`)

	req := httptest.NewRequest("POST", "/import-profile", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	handler.ImportProfile(res, req)

	if res.Code != http.StatusCreated {
		t.Errorf("Se esperaba el código de estado %d, pero se obtuvo %d", http.StatusCreated, res.Code)
	}
}

