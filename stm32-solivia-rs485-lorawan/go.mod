module github.com/ofauchon/tinygo-sandbox/stm32-solivia-rs485-lorawan

go 1.17

require (
	github.com/ofauchon/go-lorawan-stack v0.0.0-20220101213351-30d6f0ec0af3
	tinygo.org/x/drivers v0.16.0

)

require (
	github.com/TheThingsNetwork/go-cayenne-lib v1.1.0 // indirect
	github.com/jacobsa/crypto v0.0.0-20190317225127-9f44e2d11115 // indirect
	github.com/jacobsa/oglematchers v0.0.0-20150720000706-141901ea67cd // indirect
	github.com/jacobsa/oglemock v0.0.0-20150831005832-e94d794d06ff // indirect
	github.com/jacobsa/ogletest v0.0.0-20170503003838-80d50a735a11 // indirect
	github.com/jacobsa/reqtrace v0.0.0-20150505043853-245c9e0234cb // indirect
)

replace tinygo.org/x/drivers => /home/olivier/dev/contrib/tinygo-drivers

replace github.com/ofauchon/go-lorawan-stack => /home/olivier/dev/perso/go-lorawan-stack
