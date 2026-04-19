build:
	go build -o adb-install main.go

install: build
	sudo mv adb-install /usr/local/bin/

test:
	@echo "Testing with dummy APKs..."
	@touch dummy1.apk dummy2.apk
	./adb-install dummy1.apk dummy2.apk || true

clean:
	rm -f adb-install dummy*.apk

.PHONY: build install test clean
