package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/util/captcha"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12/context"
	"github.com/spf13/viper"
	"time"
)

type SessionController struct {
	Ctx         context.Context
	UserService service.UserService
}

func NewSessionController() *SessionController {
	return &SessionController{
		UserService: service.NewUserService(),
	}
}

// Login
// @Tags auth
// @Summary Login
// @Description Login
// @Param request body dto.LoginCredential true "request"
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.Profile
// @Router /auth/session/ [post]
func (s *SessionController) Post() (*dto.Profile, error) {
	aul := dto.LoginCredential{}
	if err := s.Ctx.ReadJSON(&aul); err != nil {
		return nil, err
	}
	if aul.CaptchaId != "" {
		err := captcha.VerifyCode(aul.CaptchaId, aul.Code)
		if err != nil {
			return nil, err
		}
	}
	var profile *dto.Profile
	switch aul.AuthMethod {
	case constant.AuthMethodJWT:
		p, err := s.checkSessionLogin(aul.Username, aul.Password, true)
		if err != nil {
			return nil, err
		}
		profile = p
	default:
		p, err := s.checkSessionLogin(aul.Username, aul.Password, false)
		if err != nil {
			return nil, err
		}
		profile = p
		session := constant.Sess.Start(s.Ctx)
		session.Set(constant.SessionUserKey, profile)
	}
	return profile, nil
}

// Logout
// @Tags auth
// @Summary Logout
// @Description Logout
// @Accept  json
// @Produce  json
// @Router /auth/session/ [delete]
func (s *SessionController) Delete() error {
	session := constant.Sess.Start(s.Ctx)
	session.Delete(constant.SessionUserKey)
	return nil
}

func (s *SessionController) checkSessionLogin(username string, password string, jwt bool) (*dto.Profile, error) {
	u, err := s.UserService.UserAuth(username, password)
	if err != nil {
		return nil, err
	}
	resp := &dto.Profile{}
	resp.User = toSessionUser(*u)
	if jwt {
		token, err := createToken(toSessionUser(*u))
		if err != nil {
			return nil, err
		}
		resp.Token = token
	}
	return resp, err
}

func toSessionUser(u model.User) dto.SessionUser {
	return dto.SessionUser{
		UserId:   u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Language: u.Language,
		IsActive: u.IsActive,
		IsAdmin:  u.IsAdmin,
	}
}

func createToken(user dto.SessionUser) (string, error) {
	exp := viper.GetInt("jwt.exp")
	secretKey := []byte(viper.GetString("jwt.secret"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":     user.Name,
		"email":    user.Email,
		"userId":   user.UserId,
		"isActive": user.IsActive,
		"language": user.Language,
		"isAdmin":  user.IsAdmin,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Minute * time.Duration(exp)).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}
