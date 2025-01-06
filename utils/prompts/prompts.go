package prompts

import (
	"fmt"
)

var (
	enhanceSearchPrompt = `You are a search query optimization assistant. Your task is to enhance user queries to maximize clarity, precision, and relevance for web searches. Follow these steps:

1. **Understand the Query**: Analyze the user's input query to determine its core intent, scope, and potential ambiguities.
2. **Add Context**: Incorporate relevant details such as time frame, location, format, and intended audience, as applicable.
3. **Clarify and Refine**: Resolve vagueness, specify ambiguous terms, and ensure the query has a clear purpose.
4. **Structure for Results**: Organize the query to encourage search engines to return comprehensive, high-quality results.
5. **Length Constraint**: The enhanced query MUST be less than 400 characters in length.

For instance, if the input is: "Best restaurants near me," output an enhanced query like:
"Find highly-rated restaurants within a 5-mile radius of [specific location] offering [type of cuisine] suitable for [occasion]. Include recent reviews and price ranges."

Remember: Keep the enhanced query under 400 characters while maintaining clarity and relevance.

Context: %s

Now, enhance the following query: %s`

	extractQueryPrompt = `Extract only the enhanced search query from the following text. Output should be the query only, with no quotes, markdown, or extra formatting:

%s`

	titleSystemPrompt = `Create a concise, 3-5 word phrase as a header for the following query, strictly adhering to the 3-5 word limit and avoiding the use of the word 'title', and do not generate any other text than the 3-5 word summary and do not use any markdown formatting or any asterisks for bold, and do NOT use quotation marks. 
Examples of titles:
	Stock Market Trends
	Perfect Chocolate Chip Recipe
	Evolution of Music Streaming
	Remote Work Productivity Tips
	Artificial Intelligence in Healthcare
	Video Game Development Insights`
)

// GetEnhanceSearchPrompt returns the formatted enhance search prompt
func GetEnhanceSearchPrompt(locationInfo, query string) string {
	return fmt.Sprintf(enhanceSearchPrompt, locationInfo, query)
}

// GetExtractQueryPrompt returns the formatted extract query prompt
func GetExtractQueryPrompt(enhancedQuery string) string {
	return fmt.Sprintf(extractQueryPrompt, enhancedQuery)
}

// GetTitleSystemPrompt returns the title system prompt
func GetTitleSystemPrompt() string {
	return titleSystemPrompt
} 