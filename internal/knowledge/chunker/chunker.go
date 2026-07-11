package chunker

import (
	"strings"
)

// Chunk splits a text into smaller chunks based on a maximum character limit,
// attempting to split by newlines or spaces if possible, to simulate a tokenizer limit.
func Chunk(text string, maxChars int) []string {
	if len(text) <= maxChars {
		return []string{text}
	}

	var chunks []string
	
	// A very naive implementation: split by paragraphs (double newlines), then sentences (newlines).
	// For simplicity in Phase 3, we will split by lines and accumulate until we hit maxChars.
	
	lines := strings.Split(text, "\n")
	
	var currentChunk strings.Builder
	for _, line := range lines {
		// If a single line is longer than maxChars, we just force split it
		if len(line) > maxChars {
			if currentChunk.Len() > 0 {
				chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
				currentChunk.Reset()
			}
			
			// Force split the long line
			for len(line) > maxChars {
				chunks = append(chunks, line[:maxChars])
				line = line[maxChars:]
			}
			if len(line) > 0 {
				currentChunk.WriteString(line + "\n")
			}
			continue
		}
		
		if currentChunk.Len()+len(line)+1 > maxChars {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}
		
		currentChunk.WriteString(line + "\n")
	}
	
	if currentChunk.Len() > 0 {
		finalChunk := strings.TrimSpace(currentChunk.String())
		if len(finalChunk) > 0 {
			chunks = append(chunks, finalChunk)
		}
	}

	return chunks
}
