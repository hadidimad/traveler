package main

import (
	"net/http"
)

type Route struct {
	Pattern  string
	Method   string
	Function http.HandlerFunc
}

var routes = []Route{
	{
		"/",
		"GET",
		homePageHandler,
	},
	{
		"/login",
		"GET",
		loginGetHandler,
	},
	{
		"/login",
		"POST",
		loginPostHandler,
	},
	{
		"/logout",
		"GET",
		logoutGetHandler,
	},
	{
		"/signup",
		"GET",
		signupGetHandler,
	},
	{
		"/signup",
		"POST",
		signupPostHandler,
	},
	{
		"/travels",
		"GET",
		travelsGetHandler,
	},
	{
		"/user",
		"GET",
		userGetHandler,
	},
	{
		"/useredit",
		"GET",
		userEditGetHandler,
	},
	{
		"/useredit",
		"POST",
		userEditPostHandler,
	},
	{
		"/userdelete",
		"GET",
		userDeleteGetHandler,
	},
	{
		"/travelinfo",
		"GET",
		travelInfoGetHandler,
	},
	{
		"/newtravel",
		"GET",
		newTravelGetHandler,
	},
	{
		"/newtravel",
		"POST",
		newTravelPostHandler,
	},
	{
		"/deletetravel",
		"GET",
		deleteTravelGetHandler,
	},
	{
		"/userinfo",
		"GET",
		userInfoGetHandler,
	},
	{
		"/travellike",
		"GET",
		travelLikeGetHandler,
	},
	{
		"/travelunlike",
		"GET",
		travelUnLikeGetHandler,
	},
}
