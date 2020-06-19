module github.com/cuttle-ai/websockets

go 1.13

replace github.com/cuttle-ai/auth-service => ../auth-service/

replace github.com/cuttle-ai/brain => ../brain/

replace github.com/cuttle-ai/octopus => ../octopus/

replace github.com/cuttle-ai/configs => ../configs/

replace github.com/cuttle-ai/db-toolkit => ../db-toolkit/

replace github.com/cuttle-ai/go-sdk => ../go-sdk/

require (
	github.com/cuttle-ai/auth-service v0.0.0-00010101000000-000000000000
	github.com/cuttle-ai/brain v0.0.0-00010101000000-000000000000
	github.com/cuttle-ai/configs v0.0.0-20200326184731-6eb244838d9c
	github.com/googollee/go-socket.io v1.4.3
	github.com/hashicorp/consul/api v1.4.0
	github.com/jinzhu/gorm v1.9.12
)
