init: build clear
	terraform init

clear:
	rm -rf .terraform && rm -f .terraform.lock.hcl

build :
	cd .. && make fmt && make install