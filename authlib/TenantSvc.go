package authlib

import (
	"encoding/json"

	"duov6.com/common"
	notifier "duov6.com/duonotifier/client"
	"duov6.com/gorest"
	"duov6.com/session"
	"duov6.com/term"
)

type TenantSvc struct {
	gorest.RestService
	autherized                gorest.EndPoint `method:"GET" path:"/tenant/Autherized/{TenantID:string}" output:"TenantAutherized"`
	getTenant                 gorest.EndPoint `method:"GET" path:"/tenant/GetTenant/{TenantID:string}" output:"Tenant"`
	acceptRequest             gorest.EndPoint `method:"GET" path:"/tenant/AcceptRequest/{email:string}/{RequestToken:string}" output:"bool"`
	getTenants                gorest.EndPoint `method:"GET" path:"/tenant/GetTenants/{securityToken:string}" output:"[]TenantMinimum"`
	getSampleTenantForm       gorest.EndPoint `method:"GET" path:"/tenant/GetSampleTenantForm/" output:"Tenant"`
	inviteUser                gorest.EndPoint `method:"POST" path:"/tenant/InviteUser/" postdata:"[]InviteUsers"`
	createTenant              gorest.EndPoint `method:"POST" path:"/tenant/CreateTenant/" postdata:"Tenant"`
	tenantUpgrade             gorest.EndPoint `method:"POST" path:"/tenant/TenantUpgrad/" postdata:"map[string]string"`
	searchTenants             gorest.EndPoint `method:"GET" path:"/tenant/SearchTenants/{SearchString:string}/{pagesize:int}/{startPoint:int}" output:"[]Tenant"`
	subciribe                 gorest.EndPoint `method:"GET" path:"/tenant/Subciribe/{TenantID:string}" output:"bool"`
	getUsers                  gorest.EndPoint `method:"GET" path:"/tenant/GetUsers/{TenantID:string}" output:"[]string"`
	addUser                   gorest.EndPoint `method:"GET" path:"/tenant/AddUser/{email:string}/{level:string}" output:"bool"`
	removeUser                gorest.EndPoint `method:"GET" path:"/tenant/RemoveUser/{email:string}" output:"bool"`
	tranferAdmin              gorest.EndPoint `method:"GET" path:"/tenant/TranferAdmin/{email:string}" output:"bool"`
	getPendingTenantRequest   gorest.EndPoint `method:"GET" path:"/tenant/GetPendingTenantRequest/" output:"[]PendingUserRequest"`
	getMyPendingTenantRequest gorest.EndPoint `method:"GET" path:"/tenant/GetMyPendingTenantRequest/" output:"[]PendingUserRequest"`
	getDefaultTenant          gorest.EndPoint `method:"GET" path:"/tenant/GetDefaultTenant/{UserID:string}" output:"Tenant"`
	setDefaultTenant          gorest.EndPoint `method:"GET" path:"/tenant/SetDefaultTenant/{UserID:string}/{TenantID:string}" output:"bool"`
}

func (T TenantSvc) GetDefaultTenant(UserID string) Tenant {
	//fmt.Println(T.Context.Request().Header["SecurityToken"])
	_, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	if error == "" {
		th := TenantHandler{}
		return th.GetDefaultTenant(UserID)
		//b, _ := json.Marshal(th.GetDefaultTenant(UserID))
		//T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(b)
	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
		return Tenant{}
	}
}

func (T TenantSvc) SetDefaultTenant(UserID string, TenantID string) bool {
	//fmt.Println(T.Context.Request().Header["SecurityToken"])
	_, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	if error == "" {
		th := TenantHandler{}
		return th.SetDefaultTenant(UserID, TenantID)
		//b, _ := json.Marshal(th.GetDefaultTenant(UserID))
		//T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(b)
	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
		return false
	}
}

func (T TenantSvc) CreateTenant(t Tenant) {
	//fmt.Println(T.Context.Request().Header["SecurityToken"])
	user, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	if error == "" {
		th := TenantHandler{}
		b, _ := json.Marshal(th.CreateTenant(t, user, false))
		T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(b)

	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
		return
	}
}

func (T TenantSvc) TenantUpgrade(Otherdata map[string]string) {
	user, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	if error == "" {
		th := TenantHandler{}
		t, err := th.UpgradPackage(user, Otherdata)

		if err == "" {
			b, _ := json.Marshal(t)
			T.ResponseBuilder().SetResponseCode(200).WriteAndOveride(b)
			return
		} else {
			T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(err))
			return
		}

	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
		return
	}

}

func (T TenantSvc) TranferAdmin(email string) bool {
	user, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	if error == "" {
		auth := AuthHandler{}
		u, err := auth.GetUser(email)
		if err == "" {
			th := TenantHandler{}
			return th.TransferAdmin(user, u.UserID)
		} else {
			return false
		}
	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
		return false
	}
}

func (T TenantSvc) Autherized(TenantID string) (outCrt TenantAutherized) {
	//fmt.Println(T.Context.Request().Header["SecurityToken"])
	user, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	//outCrt = TenantAutherize{}
	//TenantID
	if error == "" {
		th := TenantHandler{}
		b, d := th.Autherized(TenantID, user)

		if b {
			outCrt = d
			return d
		} else {
			T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("Tenant ID " + TenantID + " not Atherized")))
			return
		}
	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
		return
	}
}

func (T TenantSvc) GetTenant(TenantID string) Tenant {
	//fmt.Println(T.Context.Request().Header.Get("Securitytoken"))
	_, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	//outCrt = TenantAutherize{}

	if error == "" {
		th := TenantHandler{}
		return th.GetTenant(TenantID)
	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
		return Tenant{}
	}
}

func (T TenantSvc) GetUsers(TenantID string) []string {
	//fmt.Println(T.Context.Request().Header.Get("Securitytoken"))
	u, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	//outCrt = TenantAutherize{}
	//ul := []User{}
	if error == "" {
		th := TenantHandler{}
		//a := AuthHandler{}
		//a.GetUser(email)
		return th.GetUsersForTenant(u, TenantID)
	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
		return []string{}
	}
}

func (T TenantSvc) InviteUser(users []InviteUsers) {
	//fmt.Println(T.Context.Request().Header.Get("Securitytoken"))
	user, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	if error == "" {
		th := TenantHandler{}
		th.AddUserToTenant(user, users)
	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
		return
	}
}

func (T TenantSvc) AddUser(email, level string) bool {
	//fmt.Println(T.Context.Request().Header.Get("Securitytoken"))
	user, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	if error == "" {
		auth := AuthHandler{}
		a, err := auth.GetUser(email)
		if err == "" {

			th := TenantHandler{}
			t := th.GetTenant(user.Domain)
			var inputParams map[string]string
			inputParams = make(map[string]string)
			inputParams["@@EMAIL@@"] = email
			inputParams["@@INVEMAIL@@"] = user.Email
			inputParams["@@NAME@@"] = user.Name
			inputParams["@@DOMAIN@@"] = user.Domain

			go notifier.Notify("ignore", "tenant_invitation_existing", email, inputParams, nil)
			//go email.Send("ignore", "Invitation to register !", "com.duosoftware.auth", "email", "tenant_invitation", inputParams, nil, email)
			th.AddUsersToTenant(user.Domain, t.Name, a.UserID, level)
			return true

		} else {
			tmp := tempRequestGenerator{}
			th := TenantHandler{}
			t := th.GetTenant(user.Domain)
			o := make(map[string]string)
			o["process"] = "tenant_invitation"
			o["email"] = email
			o["invitedUserID"] = user.UserID
			o["name"] = user.Name
			o["domain"] = user.Domain
			o["fromuseremail"] = user.Email
			o["tname"] = t.Name
			o["level"] = level
			//o["userid"] = a.UserID
			code := tmp.GenerateRequestCode(o)
			var inputParams map[string]string
			inputParams = make(map[string]string)
			inputParams["@@EMAIL@@"] = email
			inputParams["@@INVEMAIL@@"] = user.Email
			inputParams["@@NAME@@"] = user.Name
			inputParams["@@DOMAIN@@"] = user.Domain
			inputParams["@@CODE@@"] = code

			//go notifier.Send("ignore", "User Login Notification.", "com.duosoftware.auth", "email", "tenant_invitation", inputParams, nil, email)
			go notifier.Notify("ignore", "tenant_invitation", email, inputParams, nil)
			//go email.Send("ignore", "Invitation to register !", "com.duosoftware.auth", "email", "tenant_invitation", inputParams, nil, email)
			return true
		}
	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
		return false
	}
}

func (T TenantSvc) RemoveUser(email string) bool {
	//fmt.Println(T.Context.Request().Header.Get("Securitytoken"))
	user, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	if error == "" {
		auth := AuthHandler{}
		u, err := auth.GetUser(email)
		if err == "" {
			if user.Email == email {
				T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("Cant remove yourself from the.")))
				return false
			}
			th := TenantHandler{}
			//_, p := th.Autherized(user.Domain, user)
			//if p.SecurityLevel == "admin" {
			//t := th.GetTenant(user.Domain)
			//th.AddUsersToTenant(user.Domain, t.Name, a.UserID, "level")
			//th := TenantHandler{}
			return th.RemoveUserFromTenant(u.UserID, user.Domain)
			//return true
			//} else {
			//T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte("Need to have Admin access to tenant to add user"))
			//return false
			//}

		} else {
			return false
		}
	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
		return false
	}
}

func (T TenantSvc) AcceptRequest(email, RequestToken string) bool {
	//fmt.Println(T.Context.Request().Header.Get("Securitytoken"))
	//user, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	//if error == "" {
	tmp := tempRequestGenerator{}

	o, _ := tmp.GetRequestCode(RequestToken)
	th := TenantHandler{}
	term.Write(o, term.Debug)
	term.Write(o["process"], term.Debug)
	switch o["process"] {
	case "tenant_invitation":
		auth := AuthHandler{}
		a, err := auth.GetUser(o["email"])
		if err == "" {
			th.AddUsersToTenant(o["domain"], o["tname"], a.UserID, o["level"])
			return true
		} else {
			T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("Email not registered.")))
			return false
		}
		break
	case "tenant_useradd":
		auth := AuthHandler{}
		a, err := auth.GetUser(o["email"])
		if err == "" {
			th.AddUsersToTenant(o["TenantID"], o["tname"], a.UserID, o["level"])
			th.RemovePendingRequest(o["TenantID"], a.EmailAddress)
			return true
		} else {
			T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("Email not registered.")))
			return false
		}
		break
	default:
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("Unatherized token")))
		return false
		break
	}
	//return th.AcceptRequest(user, securityLevel, RequestToken, accept)
	//} else {

	//}
	T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("Unatherized token")))
	return false
}

func (T TenantSvc) GetTenants(securityToken string) []TenantMinimum {
	tns := []TenantMinimum{}
	user, error := session.GetSession(securityToken, "Nil")

	if error == "" {
		th := TenantHandler{}
		return th.GetTenantsForUser(user.UserID)
	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
	}
	return tns
}

func (T TenantSvc) SearchTenants(SearchString string, pageSize, startPoint int) []Tenant {
	th := TenantHandler{}

	return th.SearchTenants(SearchString, pageSize, startPoint)
}

func (T TenantSvc) GetSampleTenantForm() Tenant {
	var t Tenant
	t = Tenant{}
	t.Name = "Sample Tenant"
	t.OtherData = make(map[string]string)
	t.OtherData["CompanyName"] = "DuoSoftware Pvt Ltd"
	t.OtherData["SampleAttributs"] = "Values"
	t.Private = true
	t.Statistic = make(map[string]string)
	t.Statistic["NumberOfUsers"] = "10"
	t.Statistic["DataUp"] = "1GB"
	t.Statistic["DataDown"] = "1GB"
	t.Shell = "Shell"
	t.TenantID = "smapletenat.duoworld.info"
	return t
}

func (T TenantSvc) Subciribe(TenantID string) bool {
	user, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")
	if error == "" {
		th := TenantHandler{}
		return th.RequestToTenant(user, TenantID)
	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte("SecurityToken  not Autherized"))
		return false
	}
}

type PendingUserRequest struct {
	UserID   string
	Email    string
	TenantID string
	Name     string
	Code     string
}

func (T TenantSvc) GetPendingTenantRequest() (m []PendingUserRequest) {
	var tns []PendingUserRequest
	user, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")

	if error == "" {
		th := TenantHandler{}
		tns, _ = th.GetPendingRequests(user)

	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
	}
	m = tns
	return
}

func (T TenantSvc) GetMyPendingTenantRequest() (m []PendingUserRequest) {
	var tns []PendingUserRequest
	user, error := session.GetSession(T.Context.Request().Header.Get("Securitytoken"), "Nil")

	if error == "" {
		th := TenantHandler{}
		tns, _ = th.GetMyPendingRequests(user)

	} else {
		T.ResponseBuilder().SetResponseCode(401).WriteAndOveride([]byte(common.ErrorJson("SecurityToken  not Autherized")))
	}
	m = tns
	return
}
