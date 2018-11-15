package main

import (
	"drone-k8-deploy/kube"

	//"fmt"
	"github.com/sirupsen/logrus"
)

type (
	Repo struct {
		Owner string
		Name  string
	}

	Build struct {
		Tag     string
		Event   string
		Number  int
		Commit  string
		Ref     string
		Branch  string
		Author  string
		Status  string
		Link    string
		Started int64
		Created int64
	}

	Job struct {
		Started int64
	}

	Config struct {
		Cert           string
		Server         string
		Token          string
		Namespace      string
		Template       string
		KubeConfigFile string
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Job    Job
	}
)

func (p Plugin) Exec() error {
	/*p.Config.Server = "https://172.16.0.37:6443"
	p.Config.Cert = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZuakNDQTRhZ0F3SUJBZ0lVZFRRcG8xcGVKbGM0eDJ3N3hKZ1hiY1BjK0lVd0RRWUpLb1pJaHZjTkFRRU4KQlFBd1p6RUxNQWtHQTFVRUJoTUNRMDR4RVRBUEJnTlZCQWdUQ0ZOb1pXNTZhR1Z1TVJFd0R3WURWUVFIRXdoVAphR1Z1ZW1obGJqRU1NQW9HQTFVRUNoTURhemh6TVE4d0RRWURWUVFMRXdaVGVYTjBaVzB4RXpBUkJnTlZCQU1UCkNtdDFZbVZ5Ym1WMFpYTXdIaGNOTVRnd05USXlNREV6TVRBd1doY05Nak13TlRJeE1ERXpNVEF3V2pCbk1Rc3cKQ1FZRFZRUUdFd0pEVGpFUk1BOEdBMVVFQ0JNSVUyaGxibnBvWlc0eEVUQVBCZ05WQkFjVENGTm9aVzU2YUdWdQpNUXd3Q2dZRFZRUUtFd05yT0hNeER6QU5CZ05WQkFzVEJsTjVjM1JsYlRFVE1CRUdBMVVFQXhNS2EzVmlaWEp1ClpYUmxjekNDQWlJd0RRWUpLb1pJaHZjTkFRRUJCUUFEZ2dJUEFEQ0NBZ29DZ2dJQkFMRUVBK0haRStRSm5SaXUKdncwTnNCYzBaTkVhZDlkYTVzaDA5RmlkNjhER29BSVhLK1d0a0VWWTJ4RWxubE8weTF0M2hlN2VCMnZpKzNWUQpWRnJ6VnQzNDdSRHBwcUg1Q3Avc1dlNlAyZ2tHL25paFZRSW1xS2JmdlZCSzQrMDFtUjcycnFORjFJRVlZOTNiCmwrMFUxeFlhaTJtWkpjdTJkU0RDYjVKclpIbmtSbkV1VUhxckRZYUMycjNSU1BsUnZMOWtXZnFURXN2MUZ6MUoKRnZtWG9WUHpDdWgwL21VdFd6K0lzZlZHV0NNdDdORFcwaEdjdFEvd05tRlpJcHFYQlR4UjQ5VkZ0Nmg4QVRwSgpDMmNjOGh6ODBIdUVDalNrd2VSTkFzVHFEbTBIVDR4UHhmSXVmRWE4akpwdVVRRlNnTTI1QmNvMTJVMHhkWC83CmVxOG1MNTZUSmEwLzdlcGVHcjRScW5vRlJoRnlsekRMc0N6MHl6V3c0Z0VQenk1ZjBkMXNtOUovc1c5a2x1SjUKZXRJTFA4R3NoeHh3UnczNnZkd1NQQVU5cHIzL2FpTWlYNGQrVWNOd2JKYS94QnZWOTZrYnJSbmtYS01EQU1EZwpkeDJXcDJqWnJuUkErRjMySXFJdkI5THBTOWlhRzFOMzZJUWxxOCtPQlplL1ZXcVdxeVVtY2xCRk5QUXdzMVZTCjh5WGVnRnNtMDNQbmtJVTBaOENYdHk1VEIwMzRiYURrQisyVklOQVJRWG9USVI3WkM3dnJlWDMvdGU2ckY0YloKVVdHRjhFWnRkU3pKWkM2dU1Vd1ZscndLTGhZTzh2YXRuZFBkNzFFN0liRGpub21DVEx0VnVsd0krZldGeFU4MgpPS1FwYWx0aU85N1RkT3FDQ0ZqQ0kwSTZiWXpOQWdNQkFBR2pRakJBTUE0R0ExVWREd0VCL3dRRUF3SUJCakFQCkJnTlZIUk1CQWY4RUJUQURBUUgvTUIwR0ExVWREZ1FXQkJSQytGU0hHQW4rUE9vTmxpNEs1SmRYVlJwZUtEQU4KQmdrcWhraUc5dzBCQVEwRkFBT0NBZ0VBSUo0dWhUa3ovRTY2UkRoYXJhZDdyc3hQRVRBOXoxK0lsMFQ0KzZSTgppSWJUVFQ2VVIwSTViZDd5R0tBbkxocXN5c1EzbFE2YkNkN3dHYTJMZ0hrdEhVaG90SGxVRkllaU41Sk00cm9kCkpBOWVyeFJqcjYzY09zbTBmZWp0cnUyTGgwUW84dDNSS1hQeHBNd3lxTithelZXSnhCTnRXOUJHWGVSdHE1bVIKR2RWU0h0LzJlOTlhdnpueXVOM3lCMjYvUTNKRUJwYkQ1VHF3b2o0eDRqaEc0MVdweFpjRTF1OGNET0U1RWVPMgpYRXplYlhNUE5SNmlPKyt0Q2ZtSitDTnZNRENvanF3bDJIL2lBQlcyaUN6QXliWFpxNlZDUWttMXhoR0llWkJiCjNvSzB3VUZRQkJhTjh3SUdocXVXOWkwOFdibnd0dk50U2xBNVIwOGZianBsMmFDN1czRm9IR3RaL1hjVXkrN0IKQkFYRkUvRGxTMnR5dlFSM3padWZCVFRyNE54Ty90SnZoemRNcysvZlc0bkpPWDBFVzVHcmhwRWFqaXkzMXFOQwpmYWlrWWRFTHRrY0tiNVh0Z3FMUHAyRWt3QUhZWFFCQmI4NEY2WTVsTS9EM1R3c0gwSHJqWWQ1MHBCQXpDNDdsCnl6bUlqbGQ1VG1DSzBDQ3FBWmFETllKeWdaUlVzSE9jTW9kUmlXN3gyOENvcjFUN3UzeG1FVGR0aWszNXNCZ2oKOVRKRm1kQ21ZaUhvTFdkOUR5YTNBakFGclBhd0RoTjBZYUJiK0JhTmdKVzI5TlYxYWU1Y25NV0NpNmM0SktibQpJbjg0M2xNWnJnUmZiaTMrWFBMcE4yZ2xxQlNzUEVYNzdTOGdOZ0VDTktVVHgzZStQZVN5ak1mM3hNbkNZT0FCCnhHMD0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
	p.Config.Token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJoYWt1ci1kcm9uZS1rOC1kZWxwb3ktdG9rZW4tbHFiZG0iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiaGFrdXItZHJvbmUtazgtZGVscG95Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiMjU5NTk1YzctOTBhMC0xMWU4LWE1ZDEtMGNjNDdhNTQyOGMwIiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50Omt1YmUtc3lzdGVtOmhha3VyLWRyb25lLWs4LWRlbHBveSJ9.fs4esexKT9nnT_C9OHGfFWlnhD2TuIiiVmIMjYVx-aPrKvpXRUoMod81aFylgGc1vNrFton1EV_Xf3DypuTZ8VIgGDxJRhk6S51AfSESaRipYTh_n3VdD7-AYckLo5unEHgHAVvG37UhWRKGDn4eGevsVavCiBERk-J5roYCHeKYmfBwjh1Q9OU1o5RfsJcDnfbeDFwAqemMbavez8DSqTwX3Id-MESVcpwXWXAsNLOqVTV5g5jlyO03nxU3igKGvZ0BBaK73GFjU5KclL8eCYCb-6PM0PTlB2tRVo9GoY2ghLQLwc_no4aAExMIfFVZQ6GaGiLZ9rdjdhs99y4rbyaKh2qAZUfa-7p4MG6L_wrE_zmWU2-rB_k833Rhj1HB4GRAx5_8sI5Bo-sCpSvCy7QNYTUs9ibmVLuvQd5PQVRB3RnMW_n1vrureBYCEeATSE2ydw4vYIVpcu1jN56bkVAzpLXsCVQsG7IN1YmlA4EXCM_inVjfTMv1-oByDqYxzGXWqoKLSOxwNtcHsZN3yBKLrSGV4zYFg9DzRaCJcGfWQWoQ4J0IJ982mXNrZWFPGxcjgE5krA2oqytb6MObSNxFckKUqZ_fOu4XwpkPGpCrNNRXeHiJ6j7WgLuDnbb6fQL3jIsLn46bnFOfkGezzHFeSZz8GdXCzUO_RC3kb4E"
	p.Config.Template = "rumia.yaml"*/

	/*if p.Config.Server == "" {
		logrus.Fatal("KUBERNETES_SERVER is not defined")
	}
	if p.Config.Token == "" {
		logrus.Fatal("KUBERNETES_TOKEN is not defined")
	}
	if p.Config.Cert == "" {
		logrus.Fatal("KUBERNETES_CERT is not defined")
	}
	if p.Config.Namespace == "" {
		p.Config.Namespace = "default"
	}
	if p.Config.Template == "" {
		logrus.Fatal("KUBERNETES_TEMPLATE is not defined")
	}

	err := kube.InitDefaultClient(p.Config.Server, p.Config.Cert, p.Config.Token)*/

	if p.Config.KubeConfigFile == "" {
		p.Config.KubeConfigFile = "./kubeconfig"
	}
	err := kube.InitOutClusterClient(p.Config.KubeConfigFile)
	if err != nil {
		logrus.Fatal("failed to init default kubernetes client", "->", err)
	}

	tplContent, err := RenderTemplate(p.Config.Template, p)
	if err != nil {
		logrus.Fatal("error when render template file", "->", err)
	}

	res, err := kube.ParseResources(tplContent)
	if err != nil {
		logrus.Fatal("error when parse template file", "->", err)
	}

	if err = kube.UpdateResourceList(res); err != nil {
		logrus.Fatal("failed to update resources", "->", err)
	}

	return nil
}
