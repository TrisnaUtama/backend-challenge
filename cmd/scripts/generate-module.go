package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Module name required: go run gen-module.go <module_name>")
		os.Exit(1)
	}
	module := strings.ToLower(os.Args[1])

	moduleBase := filepath.Join("internal", "modules", module)
	entitiesDir := filepath.Join("internal", "entities")

	moduleFiles := map[string]string{
		filepath.Join(moduleBase, module+".go"): fmt.Sprintf(
			"package %s\n\nimport (\n\t\"github.com/go-chi/chi/v5\"\n\t\"github.com/jackc/pgx/v5/pgxpool\"\n\n\t\"backend-challenge/internal/configs\"\n)\n\nfunc Init(r chi.Router, db *pgxpool.Pool, cfg *configs.Setting) {\n\trepo    := NewRepository(db)\n\tservice := NewService(repo, cfg)\n\thandler := NewHandler(service)\n\n\tr.Route(\"/%s\", func(r chi.Router) {\n\t\t_ = handler\n\t})\n}\n",
			module, module,
		),
		filepath.Join(moduleBase, "dto.go"): fmt.Sprintf(
			"package %s\n\n// Request & Response DTOs\n", module,
		),
		filepath.Join(moduleBase, "handler.go"): fmt.Sprintf(
			"package %s\n\nimport \"net/http\"\n\ntype Handler struct {\n\tservice Service\n}\n\nfunc NewHandler(service Service) *Handler {\n\treturn &Handler{service: service}\n}\n\n// Example:\n// func (h *Handler) List(w http.ResponseWriter, r *http.Request) {}\n\n var _ *http.Request // prevent unused import\n",
			module,
		),
		filepath.Join(moduleBase, "repository.go"): fmt.Sprintf(
			"package %s\n\nimport \"github.com/jackc/pgx/v5/pgxpool\"\n\ntype Repository interface{}\n\ntype repository struct {\n\tdb *pgxpool.Pool\n}\n\nfunc NewRepository(db *pgxpool.Pool) Repository {\n\treturn &repository{db: db}\n}\n",
			module,
		),
		filepath.Join(moduleBase, "service.go"): fmt.Sprintf(
			"package %s\n\nimport \"backend-challenge/internal/configs\"\n\ntype Service interface{}\n\ntype service struct {\n\trepo Repository\n\tcfg  *configs.Setting\n}\n\nfunc NewService(repo Repository, cfg *configs.Setting) Service {\n\treturn &service{repo: repo, cfg: cfg}\n}\n",
			module,
		),
	}

	entityFile := filepath.Join(entitiesDir, module+".go")
	entityContent := fmt.Sprintf(
		"package entities\n\nimport \"time\"\n\ntype %s struct {\n\tID        string\n\tCreatedAt time.Time\n\tUpdatedAt time.Time\n\tDeletedAt *time.Time\n}\n",
		strings.Title(module),
	)

	if err := os.MkdirAll(moduleBase, os.ModePerm); err != nil {
		fmt.Printf("Failed to create module dir: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(entitiesDir, os.ModePerm); err != nil {
		fmt.Printf("Failed to create entities dir: %v\n", err)
		os.Exit(1)
	}

	for path, content := range moduleFiles {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				fmt.Printf("Failed to write %s: %v\n", path, err)
				os.Exit(1)
			}
			fmt.Printf("✅ Created: %s\n", path)
		} else {
			fmt.Printf("⏭️  Skipped (exists): %s\n", path)
		}
	}

	if _, err := os.Stat(entityFile); os.IsNotExist(err) {
		if err := os.WriteFile(entityFile, []byte(entityContent), 0644); err != nil {
			fmt.Printf("Failed to write entity %s: %v\n", entityFile, err)
			os.Exit(1)
		}
		fmt.Printf("✅ Created entity: %s\n", entityFile)
	} else {
		fmt.Printf("⏭️  Skipped (exists): %s\n", entityFile)
	}

	fmt.Printf("\n🚀 Module '%s' created successfully!\n", module)
	fmt.Printf("📁 Module  : %s\n", moduleBase)
	fmt.Printf("📦 Entity  : %s\n", entityFile)
}
