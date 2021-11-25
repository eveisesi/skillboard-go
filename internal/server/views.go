package server

// func (s *server) handleRenderCharacterPage(w http.ResponseWriter, r *http.Request) {
// 	t, err := template.New("character.gohtml").ParseGlob("internal/server/views/*.html")
// 	if err != nil {
// 		s.writeError(r.Context(), w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	w.Header().Set(headers.ContentType, "text/html")
// 	err = t.Execute(w, struct{ Name string }{Name: "David"})
// 	if err != nil {
// 		s.writeError(r.Context(), w, http.StatusInternalServerError, err)
// 		return
// 	}
// }
