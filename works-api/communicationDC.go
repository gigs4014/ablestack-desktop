package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

//login
func postLogin(id string, password string) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")
	params := url.Values{
		"username": {id},
		"password": {password},
	}
	client := http.Client{}
	resp, err := client.PostForm(DCInfo+"/v1/login", params)
	//if err != nil {
	//	log.Error(err)
	//	log.Errorf("%v\n", resp)
	//} else {
	//	json.NewDecoder(resp.Body).Decode(&res)
	//}

	return resp, err
}

//************************** User 관련 시작 **************************
//유저를 DC 에 추가하는 func
func postDCUser(user User) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")
	params := url.Values{
		"username":   {user.Cn},
		"password":   {user.Password},
		"sn":         {user.Sn},
		"givenName":  {user.GivenName},
		"email":      {user.Mail},
		"phone":      {user.TelephoneNumber},
		"title":      {user.Title},
		"department": {user.Department},
	}
	log.Infof("paramsInfo = [%v]", params)
	log.Infof("user = [%v]", user)
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.PostForm(DCInfo+"/v1/user", params)

	return resp, err
}

// 유저 리스트를 DC 에서 조회하는 func
func getUserList() ([]User, error) {
	var DCInfo = os.Getenv("DCUrl")
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	//resp, err := client.PostForm(DCInfo+"/v1/user/", params)
	resp, err := client.Get(DCInfo + "/v1/user/")
	var res []map[string]interface{}
	var userInfoList []User
	if err != nil {
		log.WithFields(logrus.Fields{
			"authDAO": "getUserList",
		}).Errorf("DC communication failed. [%v], error [%v]", resp, err)
	} else {
		err := json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			return nil, err
		}

		jsonMarshal, _ := json.Marshal(res)
		json.Unmarshal([]byte(jsonMarshal), &userInfoList)
		log.WithFields(logrus.Fields{
			"authDAO": "getUserList",
		}).Infof("DC communication result. resp.StatusCode [%v], userInfoList [%v]", resp.StatusCode, userInfoList)
		//json.NewDecoder(resp.Body).Decode(&res)
	}

	return userInfoList, err
}

// 유저 정보를 DC 에서 조회하는 func
func getUserInfo(id string) *http.Response {
	var DCInfo = os.Getenv("DCUrl")
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	//resp, err := client.PostForm(DCInfo+"/v1/user/", params)
	resp, err := client.Get(DCInfo + "/v1/user/" + id)
	if err != nil {
		log.Error(err)
		log.Error(resp)
	}

	return resp
}

func deleteDCUser(userName string) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/v1/user/%v", DCInfo, userName), nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("유저 삭제중 에러가 발생했습니다. [%v]", err)
	}
	return resp, err
}

//************************** User 관련 끝 **************************

func insertGroup(groupName string) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")
	params := url.Values{
		"groupname": {groupName},
	}
	log.Infof("paramsInfo = [%v]", params)
	client := http.Client{}
	resp, err := client.PostForm(DCInfo+"/v1/group", params)

	log.Infof("%v %v", resp, err)
	return resp, err
}

func insertPolicyRemotefx(groupName string) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")
	client := http.Client{
		Timeout: 20 * time.Second,
	}
	params := url.Values{
		"groupname": {groupName},
	}
	log.Infof("paramsInfo = [%v]", params)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/v1/policy/remotefx/%v", DCInfo, groupName), nil)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	log.WithFields(logrus.Fields{
		"dcDAO": "insertPolicyRemotefx",
	}).Infof("workspace group add remotrfx status. resp [%v], err [%v]", resp, err)
	return resp, err
}

func insertPolicyDC(groupName string, policy Policy) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")
	client := http.Client{
		Timeout: 20 * time.Second,
	}
	params := url.Values{
		"groupname": {groupName},
	}
	log.Infof("paramsInfo = [%v]", params)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/v1/policy/%v/%v", DCInfo, policy.AdGroupName, groupName), nil)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	log.WithFields(logrus.Fields{
		"dcDAO": "insertPolicyRemotefx",
	}).Infof("workspace group add remotrfx status. resp [%v], err [%v]", resp, err)
	return resp, err
}

func deletePolicyDC(groupName string, policy Policy) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")
	client := http.Client{
		Timeout: 20 * time.Second,
	}
	params := url.Values{
		"groupname": {groupName},
	}
	log.Infof("paramsInfo = [%v]", params)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/v1/policy/%v/%v", DCInfo, policy.AdGroupName, groupName), nil)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	log.WithFields(logrus.Fields{
		"dcDAO": "insertPolicyRemotefx",
	}).Infof("workspace group add remotrfx status. resp [%v], err [%v]", resp, err)
	return resp, err
}

func selectDCWorkspacePolicyList(workspace Workspace) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(DCInfo + "/v1/group/" + workspace.Name + "/policy")

	return resp, err
}

func selectGroupList() (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(DCInfo + "/v1/group/")

	return resp, err
}

func selectGroupDetail(groupName string) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")

	client := http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Get(DCInfo + "/v1/group/" + groupName)

	return resp, err
}

func deleteGroup(groupName string) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	//resp, err := client.Get(DCInfo + "/v1/group/" + groupName)
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/v1/group/%v", DCInfo, groupName), nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	log.WithFields(logrus.Fields{
		"communicationDC": "deleteGroup",
	}).Infof("workspace group delete status. resp [%v], err [%v]", resp, err)

	return resp, err
}

func deleteComputer(computerName string) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	//resp, err := client.Get(DCInfo + "/v1/group/" + groupName)
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/v1/computer/%v", DCInfo, computerName), nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	log.WithFields(logrus.Fields{
		"authDAO": "deleteComputer",
	}).Infof("workspace group delete status. resp [%v], err [%v]", resp, err)

	return resp, err
}

func insertAddUserToGroup(groupName string, userName string) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/v1/group/%v/%v", DCInfo, groupName, userName), nil)
	if err != nil {
		log.Errorf("유저 삭제중 에러가 발생했습니다. [%v]", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err1 := client.Do(req)
	return resp, err1
}

func deleteAddUserToGroup(groupName string, userName string) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/v1/group/%v/%v", DCInfo, groupName, userName), nil)
	if err != nil {
		log.Errorf("유저 삭제중 에러가 발생했습니다. [%v]", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err1 := client.Do(req)
	return resp, err1
}

func insertUserAllocatedInstance1(username string, instanceUuid string) *http.Response {
	var DCInfo = os.Getenv("DCUrl")
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resultInstance, _ := selectInstanceList(instanceUuid, InstanceString)
	instanceInfo := resultInstance[0]
	resultUser := selectUserDBDetail(username)
	paramsMold := []MoldParams{
		{"id": instanceInfo.MoldUuid},
	}
	resultMoldInstance := getListVirtualMachinesMetrics(paramsMold)
	listVirtualMachinesMetrics := ListVirtualMachinesMetrics{}
	virtualMachineInfo, _ := json.Marshal(resultMoldInstance["listvirtualmachinesmetricsresponse"])
	json.Unmarshal([]byte(virtualMachineInfo), &listVirtualMachinesMetrics)
	log.Infof("%v", listVirtualMachinesMetrics)
	params := url.Values{
		"connection": {instanceInfo.Name},
		"parameter":  {"hostname=" + listVirtualMachinesMetrics.Virtualmachine[0].Nic[0].Ipaddress + ",port=3389,ignore-cert=true,username=" + username + ",password=" + resultUser.Password + ",domain=testdomain"},
	}
	resp, err := client.PostForm(DCInfo+"/v1/user/"+username, params)
	//req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/v1/user/%v", DCInfo, username), params)
	if err != nil {
		log.Errorf("유저 삭제중 에러가 발생했습니다. [%v]", err)
	}

	log.Infof(resp.Status)
	return resp
}

func insertConnection(username string, connectName string, parameter string) *http.Response {
	var DCInfo = os.Getenv("DCUrl")
	client := http.Client{}

	params := url.Values{
		"username":  {username},
		"parameter": {parameter},
	}
	//resp, err := client.PostForm(DCInfo+"/v1/connection/"+connectName, params)
	resp, err := client.PostForm(fmt.Sprintf("%v/v1/connection/%v", DCInfo, connectName), params)
	if err != nil {
		log.Errorf("An error occurred while creating the connection information. [%v]", err)
	}

	log.Infof(resp.Status)
	return resp
}

func delConnection(connectName string) *http.Response {
	var DCInfo = os.Getenv("DCUrl")
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/v1/connection/%v", DCInfo, connectName), nil)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(logrus.Fields{
			"communicationMold.go": "delConnection",
		}).Errorf("An error occurred while deleting the connection. resp [%v], err [%v]", resp, err)
	}

	return resp
}

func selectPasswordConvert(password string) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	aaa := DCInfo + "/v1/cmd/?cmd="
	bbb := url.QueryEscape("(\"" + password + "\" | ConvertTo-SecureString -AsPlainText -Force) | ConvertFrom-SecureString")
	resp, err := client.Get(aaa + bbb + "&timeout=20")

	log.Infof("resp [%v], err [%v]", resp, err)

	return resp, err
}

func patchUserPassword(userName string, newPassword string) (*http.Response, error) {
	var DCInfo = os.Getenv("DCUrl")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	params := url.Values{
		"password": {newPassword},
	}
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%v/v1/user/%v", DCInfo, userName), strings.NewReader(params.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(logrus.Fields{
			"communicationMold.go": "patchUserPassword",
		}).Errorf("An error occurred while deleting the connection. resp [%v], err [%v]", resp, err)
	}

	return resp, err
}
