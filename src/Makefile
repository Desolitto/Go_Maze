# Определяем переменные
BINARY_NAME = maze
SRC_DIR = .
BUILD_DIR = build
INSTALL_DIR = ../$(BUILD_DIR)  # Укажите желаемый путь установки
DIST_DIR =  ../dist  # Папка для дистрибутива
ARCHIVE_NAME = maze_dist.tar  # Имя архива

# Цели
.PHONY: all install uninstall clean dvi dist tests

# Основная цель
all: clean build

# Сборка
build:
	@echo "Сборка проекта..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)/cmd/maze/maze.go

# Установка
install: all
	@echo "Установка в $(INSTALL_DIR)..."
	mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)

# Удаление
uninstall:
	@echo "Удаление из $(INSTALL_DIR)..."
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)

# Очистка
clean:
	@echo "Очистка..."
	rm -rf $(BUILD_DIR) $(DIST_DIR)

# Генерация документации (если есть)
dvi:
	@echo "Генерация документации..."

# Распределение (например, создание архива)
dist:
	@echo "Создание дистрибутива..."
	@mkdir -p $(DIST_DIR)
	@go build -o $(BINARY_NAME) $(SRC_DIR)/cmd/maze/maze.go
	@cp ../README_RUS.md .
	@tar -cvf $(ARCHIVE_NAME) $(BINARY_NAME) README_RUS.md
	@mv $(ARCHIVE_NAME) $(DIST_DIR)
	@rm README_RUS.md $(BINARY_NAME)
	@echo "Дистрибутив создан."

# Тесты
tests:
	@echo "Запуск тестов..."
	go test ./...
