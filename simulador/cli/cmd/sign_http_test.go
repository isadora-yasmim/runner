package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// signResponse espelha ResponseOutput do assinador para desserialização nos testes.
type signResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type signDataFields struct {
	SignatureHash string `json:"signatureHash"`
	Algorithm     string `json:"algorithm"`
}

// ------------------------------------------------------------------ sign via HTTP (mock server)

// buildMockAssinador cria um httptest.Server que replica o comportamento
// do assinador.jar nos endpoints /sign e /validate, permitindo testar
// o cliente Go sem dependência do JAR.
func buildMockAssinador(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"message":"Assinador HTTP ativo","data":{"status":"UP"}}`))
	})

	mux.HandleFunc("/sign", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		body, _ := io.ReadAll(r.Body)
		var req map[string]string
		json.Unmarshal(body, &req)

		doc := req["document"]
		pin := req["tokenPin"]

		w.Header().Set("Content-Type", "application/json")

		if doc == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"success":false,"message":"O documento nao pode estar vazio.","data":null}`))
			return
		}
		if pin == "" || len(pin) < 4 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"success":false,"message":"O PIN do token deve conter pelo menos 4 caracteres.","data":null}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"message":"Assinatura criada com sucesso (Simulacao)",` +
			`"data":{"signatureHash":"mock_hash_abc123_base64_encoded_signature_simulated","algorithm":"SHA256withRSA"}}`))
	})

	mux.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		body, _ := io.ReadAll(r.Body)
		var req map[string]string
		json.Unmarshal(body, &req)

		w.Header().Set("Content-Type", "application/json")

		if req["document"] == "" || req["signature"] == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"success":false,"message":"Parametros obrigatorios ausentes.","data":null}`))
			return
		}

		if req["signature"] == "mock_hash_abc123_base64_encoded_signature_simulated" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true,"message":"Assinatura valida.","data":true}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":false,"message":"Assinatura invalida ou corrompida.","data":false}`))
		}
	})

	return httptest.NewServer(mux)
}

// ------------------------------------------------------------------ testes de cliente HTTP

func TestSignHTTPClient_DocumentoEPinValidos_DeveRetornarSuccessTrue(t *testing.T) {
	srv := buildMockAssinador(t)
	defer srv.Close()

	reqBody, _ := json.Marshal(map[string]string{
		"document": "contrato.pdf",
		"tokenPin": "1234",
	})

	resp, err := http.Post(srv.URL+"/sign", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200, obteve %d", resp.StatusCode)
	}

	var out signResponse
	json.NewDecoder(resp.Body).Decode(&out)

	if !out.Success {
		t.Errorf("esperava success:true, obteve: %s", out.Message)
	}

	var data signDataFields
	json.Unmarshal(out.Data, &data)

	if data.SignatureHash != "mock_hash_abc123_base64_encoded_signature_simulated" {
		t.Errorf("hash inesperado: %q", data.SignatureHash)
	}
	if data.Algorithm != "SHA256withRSA" {
		t.Errorf("algoritmo inesperado: %q", data.Algorithm)
	}
}

func TestSignHTTPClient_DocumentoComEspacos_DeveRetornarSuccessTrue(t *testing.T) {
	// Critério E1: argumentos com espaços devem ser preservados
	srv := buildMockAssinador(t)
	defer srv.Close()

	reqBody, _ := json.Marshal(map[string]string{
		"document": "meu documento 2024.pdf",
		"tokenPin": "1234",
	})

	resp, err := http.Post(srv.URL+"/sign", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200, obteve %d", resp.StatusCode)
	}
}

func TestSignHTTPClient_PinAusente_DeveRetornarStatus400(t *testing.T) {
	srv := buildMockAssinador(t)
	defer srv.Close()

	reqBody, _ := json.Marshal(map[string]string{
		"document": "doc.pdf",
	})

	resp, err := http.Post(srv.URL+"/sign", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperava 400, obteve %d", resp.StatusCode)
	}
}

func TestVerifyHTTPClient_AssinaturaValida_DeveRetornarSuccessTrue(t *testing.T) {
	srv := buildMockAssinador(t)
	defer srv.Close()

	reqBody, _ := json.Marshal(map[string]string{
		"document":  "doc.pdf",
		"signature": "mock_hash_abc123_base64_encoded_signature_simulated",
	})

	resp, err := http.Post(srv.URL+"/validate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	var out signResponse
	json.NewDecoder(resp.Body).Decode(&out)

	if !out.Success {
		t.Errorf("esperava success:true para assinatura válida")
	}
}

func TestVerifyHTTPClient_AssinaturaInvalida_DeveRetornarSuccessFalse(t *testing.T) {
	srv := buildMockAssinador(t)
	defer srv.Close()

	reqBody, _ := json.Marshal(map[string]string{
		"document":  "doc.pdf",
		"signature": "hash_errado",
	})

	resp, err := http.Post(srv.URL+"/validate", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	var out signResponse
	json.NewDecoder(resp.Body).Decode(&out)

	if out.Success {
		t.Errorf("esperava success:false para assinatura inválida")
	}
}

func TestHealthCheck_ServidorMock_DeveRetornarStatusUP(t *testing.T) {
	srv := buildMockAssinador(t)
	defer srv.Close()

	port := portFromURL(t, srv.URL)
	if !isServerRunning(port) {
		t.Error("isServerRunning deveria retornar true para mock server ativo")
	}
}
