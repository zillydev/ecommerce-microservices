APPS = user-service notification-service product-service order-service gateway-service

build:
	@mkdir -p bin
	@for app in $(APPS); do \
		echo "Building $$app..."; \
		go build -o bin/$$app.exe ./cmd/$$app; \
	done

run:
	@for app in $(APPS); do \
		echo "Running $$app..."; \
		bin/$$app.exe & \
	done
	@wait

clean:
	@rm -rf bin/