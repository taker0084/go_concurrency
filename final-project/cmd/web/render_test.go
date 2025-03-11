package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfig_AddDefaultData(t *testing.T) {
	req, _ := http.NewRequest("Get", "/", nil)

	ctx := getCtx(req)

	req = req.WithContext(ctx)

	testApp.Session.Put(ctx, "flash", "flash")
	testApp.Session.Put(ctx, "warning", "warning")
	testApp.Session.Put(ctx, "error", "error")

	td := testApp.addDefaultData(&TemplateData{}, req)

	if td.Flash != "flash" {
		t.Errorf("Expected flash to be flash, but got %s", td.Flash)
	}
	if td.Warning != "warning" {
		t.Errorf("Expected warning to be warning, but got %s", td.Warning)
	}
	if td.Error != "error" {
		t.Errorf("Expected error to be error, but got %s", td.Error)
	}
}

func TestConfig_IsAuthenticated(t *testing.T) {
	req, _ := http.NewRequest("Get", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	isAuth := testApp.IsAuthenticated(req)

	if isAuth {
		t.Errorf("Expected isAuth to be false, but got %v", isAuth)
	}

	testApp.Session.Put(ctx, "userID", 1)

	isAuth = testApp.IsAuthenticated(req)

	if !isAuth {
		t.Errorf("Expected isAuth to be true, but got %v", isAuth)
	}
}

func TestConfig_render(t *testing.T) {
	pathToTemplates = "./templates"

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("Get", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	testApp.render(rr, req, "home.page.gohtml", &TemplateData{})

	if rr.Code != 200 {
		t.Error("failed to render page")
	}
}
