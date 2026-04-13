package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strings"
)

type ClassifyRequest struct {
    URL string `json:"url"`
}

type ClassifyResponse struct {
    Success        bool     `json:"success"`
    Data           *Result  `json:"data,omitempty"`
    Error          string   `json:"error,omitempty"`
}

type Result struct {
    DocumentType     string   `json:"document_type"`
    DocumentTypeRU   string   `json:"document_type_ru"`
    Confidence       float64  `json:"confidence"`
    DetectedKeywords []string `json:"detected_keywords"`
    TextPreview      string   `json:"text_preview,omitempty"`
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    mux := http.NewServeMux()
    mux.HandleFunc("POST /classify", classifyHandler)
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    })

    log.Printf("📄 Document Classifier API запущен на порту %s", port)
    log.Printf("📡 POST /classify - определение типа документа")
    log.Printf("❤️  GET /health")

    http.ListenAndServe(":"+port, mux)
}

func classifyHandler(w http.ResponseWriter, r *http.Request) {
    var req ClassifyRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, "Invalid JSON", 400)
        return
    }

    if req.URL == "" {
        writeError(w, "url is required", 400)
        return
    }

    // Демо-версия: определяем тип документа по URL или заглушке
    // В реальном API здесь будет загрузка файла и анализ текста
    result := classifyByURL(req.URL)

    writeJSON(w, 200, ClassifyResponse{
        Success: true,
        Data:    result,
    })
}

func classifyByURL(url string) *Result {
    urlLower := strings.ToLower(url)

    // Простая демо-логика для MVP
    switch {
    case strings.Contains(urlLower, "passport") || strings.Contains(urlLower, "pasport"):
        return &Result{
            DocumentType:     "passport",
            DocumentTypeRU:   "Паспорт",
            Confidence:       0.85,
            DetectedKeywords: []string{"паспорт", "серия", "номер"},
            TextPreview:      "ПАСПОРТ... Серия 1234 Номер 567890...",
        }
    case strings.Contains(urlLower, "contract") || strings.Contains(urlLower, "dogovor"):
        return &Result{
            DocumentType:     "contract",
            DocumentTypeRU:   "Договор",
            Confidence:       0.80,
            DetectedKeywords: []string{"договор", "стороны", "предмет договора"},
            TextPreview:      "ДОГОВОР №123...",
        }
    case strings.Contains(urlLower, "receipt") || strings.Contains(urlLower, "chek"):
        return &Result{
            DocumentType:     "receipt",
            DocumentTypeRU:   "Чек",
            Confidence:       0.75,
            DetectedKeywords: []string{"чек", "итого", "спасибо"},
            TextPreview:      "ЧЕК... Товары: 3 шт... Итого: 1500 руб...",
        }
    case strings.Contains(urlLower, "invoice") || strings.Contains(urlLower, "schet"):
        return &Result{
            DocumentType:     "invoice",
            DocumentTypeRU:   "Счёт",
            Confidence:       0.80,
            DetectedKeywords: []string{"счёт", "инвойс", "к оплате"},
            TextPreview:      "СЧЁТ №456 от 01.01.2024...",
        }
    case strings.Contains(urlLower, "resume") || strings.Contains(urlLower, "cv") || strings.Contains(urlLower, "rezume"):
        return &Result{
            DocumentType:     "resume",
            DocumentTypeRU:   "Резюме",
            Confidence:       0.85,
            DetectedKeywords: []string{"образование", "опыт работы", "навыки"},
            TextPreview:      "ИВАНОВ ИВАН... Образование: МГУ... Опыт работы: 5 лет...",
        }
    default:
        return &Result{
            DocumentType:     "unknown",
            DocumentTypeRU:   "Неизвестно",
            Confidence:       0.50,
            DetectedKeywords: []string{},
            TextPreview:      "Не удалось определить тип документа",
        }
    }
}

func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(ClassifyResponse{
        Success: false,
        Error:   message,
    })
}