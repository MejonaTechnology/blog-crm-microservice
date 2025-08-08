package seo

import (
	"crypto/md5"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SEOAnalyzer analyzes content for SEO optimization opportunities
type SEOAnalyzer struct{}

// NewSEOAnalyzer creates a new SEO analyzer
func NewSEOAnalyzer() *SEOAnalyzer {
	return &SEOAnalyzer{}
}

// AnalyzeContent performs comprehensive SEO analysis on blog content
func (sa *SEOAnalyzer) AnalyzeContent(content ContentData) SEOAnalysis {
	analysis := SEOAnalysis{
		ContentID:   content.ID,
		Title:       content.Title,
		URL:         content.URL,
		AnalyzedAt:  time.Now(),
	}

	// Analyze title optimization
	analysis.TitleAnalysis = sa.analyzeTitleSEO(content.Title, content.PrimaryKeyword)

	// Analyze meta description
	analysis.MetaAnalysis = sa.analyzeMetaDescription(content.MetaDescription, content.PrimaryKeyword)

	// Analyze content structure and headings
	analysis.StructureAnalysis = sa.analyzeContentStructure(content.Content, content.Headings)

	// Analyze keyword optimization
	analysis.KeywordAnalysis = sa.analyzeKeywordOptimization(content)

	// Analyze readability
	analysis.ReadabilityAnalysis = sa.analyzeReadability(content.Content)

	// Analyze technical SEO factors
	analysis.TechnicalAnalysis = sa.analyzeTechnicalSEO(content)

	// Analyze link structure
	analysis.LinkAnalysis = sa.analyzeLinkStructure(content.InternalLinks, content.ExternalLinks)

	// Analyze image optimization
	analysis.ImageAnalysis = sa.analyzeImageOptimization(content.Images)

	// Calculate overall SEO score
	analysis.OverallScore = sa.calculateOverallSEOScore(analysis)

	// Generate recommendations
	analysis.Recommendations = sa.generateRecommendations(analysis)

	// Identify optimization opportunities
	analysis.Opportunities = sa.identifyOpportunities(analysis)

	return analysis
}

// analyzeTitleSEO analyzes title tag optimization
func (sa *SEOAnalyzer) analyzeTitleSEO(title, primaryKeyword string) TitleAnalysis {
	analysis := TitleAnalysis{
		Title:  title,
		Length: len(title),
	}

	// Check title length (optimal: 50-60 characters)
	if analysis.Length >= 50 && analysis.Length <= 60 {
		analysis.LengthScore = 100
		analysis.LengthStatus = "optimal"
	} else if analysis.Length >= 40 && analysis.Length < 70 {
		analysis.LengthScore = 80
		analysis.LengthStatus = "good"
	} else if analysis.Length < 40 {
		analysis.LengthScore = 60
		analysis.LengthStatus = "too_short"
	} else {
		analysis.LengthScore = 40
		analysis.LengthStatus = "too_long"
	}

	// Check primary keyword presence
	titleLower := strings.ToLower(title)
	keywordLower := strings.ToLower(primaryKeyword)
	
	if keywordLower != "" {
		if strings.Contains(titleLower, keywordLower) {
			analysis.ContainsPrimaryKeyword = true
			
			// Check keyword position (beginning is better)
			position := strings.Index(titleLower, keywordLower)
			if position == 0 {
				analysis.KeywordPosition = "beginning"
				analysis.KeywordScore = 100
			} else if position <= len(title)/2 {
				analysis.KeywordPosition = "middle"
				analysis.KeywordScore = 80
			} else {
				analysis.KeywordPosition = "end"
				analysis.KeywordScore = 60
			}
		} else {
			analysis.ContainsPrimaryKeyword = false
			analysis.KeywordScore = 20
		}
	} else {
		analysis.KeywordScore = 50 // Neutral if no primary keyword defined
	}

	// Check for power words
	powerWords := []string{"ultimate", "complete", "guide", "best", "top", "how", "why", "what", "when", "expert", "proven", "essential", "amazing", "incredible", "powerful"}
	for _, word := range powerWords {
		if strings.Contains(titleLower, word) {
			analysis.PowerWords = append(analysis.PowerWords, word)
		}
	}
	analysis.PowerWordScore = math.Min(float64(len(analysis.PowerWords))*20, 100)

	return analysis
}

// analyzeMetaDescription analyzes meta description optimization
func (sa *SEOAnalyzer) analyzeMetaDescription(metaDescription, primaryKeyword string) MetaAnalysis {
	analysis := MetaAnalysis{
		MetaDescription: metaDescription,
		Length:          len(metaDescription),
	}

	// Check meta description length (optimal: 150-160 characters)
	if analysis.Length >= 150 && analysis.Length <= 160 {
		analysis.LengthScore = 100
		analysis.LengthStatus = "optimal"
	} else if analysis.Length >= 120 && analysis.Length <= 170 {
		analysis.LengthScore = 80
		analysis.LengthStatus = "good"
	} else if analysis.Length < 120 {
		analysis.LengthScore = 60
		analysis.LengthStatus = "too_short"
	} else {
		analysis.LengthScore = 40
		analysis.LengthStatus = "too_long"
	}

	// Check primary keyword presence
	if primaryKeyword != "" {
		metaLower := strings.ToLower(metaDescription)
		keywordLower := strings.ToLower(primaryKeyword)
		
		if strings.Contains(metaLower, keywordLower) {
			analysis.ContainsPrimaryKeyword = true
			analysis.KeywordScore = 100
		} else {
			analysis.ContainsPrimaryKeyword = false
			analysis.KeywordScore = 20
		}
	} else {
		analysis.KeywordScore = 50
	}

	// Check for call-to-action words
	ctaWords := []string{"learn", "discover", "find out", "get", "download", "read", "explore", "try", "start", "join"}
	metaLower := strings.ToLower(metaDescription)
	for _, cta := range ctaWords {
		if strings.Contains(metaLower, cta) {
			analysis.CallToAction = true
			break
		}
	}
	
	if analysis.CallToAction {
		analysis.CTAScore = 100
	} else {
		analysis.CTAScore = 30
	}

	return analysis
}

// analyzeContentStructure analyzes heading structure and content organization
func (sa *SEOAnalyzer) analyzeContentStructure(content string, headings []HeadingData) StructureAnalysis {
	analysis := StructureAnalysis{
		WordCount: sa.countWords(content),
		Headings:  headings,
	}

	// Analyze word count (optimal: 1000-3000 words for blog posts)
	if analysis.WordCount >= 1000 && analysis.WordCount <= 3000 {
		analysis.WordCountScore = 100
		analysis.WordCountStatus = "optimal"
	} else if analysis.WordCount >= 500 && analysis.WordCount < 4000 {
		analysis.WordCountScore = 80
		analysis.WordCountStatus = "good"
	} else if analysis.WordCount < 500 {
		analysis.WordCountScore = 50
		analysis.WordCountStatus = "too_short"
	} else {
		analysis.WordCountScore = 70
		analysis.WordCountStatus = "very_long"
	}

	// Analyze heading structure
	h1Count := 0
	h2Count := 0
	h3Count := 0
	
	for _, heading := range headings {
		switch heading.Level {
		case 1:
			h1Count++
		case 2:
			h2Count++
		case 3:
			h3Count++
		}
	}

	analysis.H1Count = h1Count
	analysis.H2Count = h2Count
	analysis.H3Count = h3Count

	// Check H1 optimization (should have exactly 1)
	if h1Count == 1 {
		analysis.H1Score = 100
		analysis.H1Status = "optimal"
	} else if h1Count == 0 {
		analysis.H1Score = 20
		analysis.H1Status = "missing"
	} else {
		analysis.H1Score = 60
		analysis.H1Status = "multiple"
	}

	// Check H2 structure (should have at least 2-3 for good structure)
	if h2Count >= 2 && h2Count <= 8 {
		analysis.H2Score = 100
		analysis.H2Status = "good"
	} else if h2Count == 1 {
		analysis.H2Score = 70
		analysis.H2Status = "minimal"
	} else if h2Count == 0 {
		analysis.H2Score = 40
		analysis.H2Status = "missing"
	} else {
		analysis.H2Score = 80
		analysis.H2Status = "many"
	}

	// Calculate content structure score
	paragraphCount := strings.Count(content, "\n\n") + 1
	if analysis.WordCount > 0 {
		analysis.AvgWordsPerParagraph = float64(analysis.WordCount) / float64(paragraphCount)
	}

	// Optimal paragraph length: 50-100 words
	if analysis.AvgWordsPerParagraph >= 50 && analysis.AvgWordsPerParagraph <= 100 {
		analysis.ParagraphScore = 100
	} else if analysis.AvgWordsPerParagraph >= 30 && analysis.AvgWordsPerParagraph <= 150 {
		analysis.ParagraphScore = 80
	} else {
		analysis.ParagraphScore = 60
	}

	return analysis
}

// analyzeKeywordOptimization analyzes keyword usage and density
func (sa *SEOAnalyzer) analyzeKeywordOptimization(content ContentData) KeywordAnalysis {
	analysis := KeywordAnalysis{
		PrimaryKeyword:    content.PrimaryKeyword,
		SecondaryKeywords: content.SecondaryKeywords,
	}

	contentText := content.Content
	wordCount := sa.countWords(contentText)
	
	if wordCount == 0 {
		return analysis
	}

	contentLower := strings.ToLower(contentText)

	// Analyze primary keyword
	if content.PrimaryKeyword != "" {
		primaryKeywordLower := strings.ToLower(content.PrimaryKeyword)
		
		// Count primary keyword occurrences
		analysis.PrimaryKeywordCount = strings.Count(contentLower, primaryKeywordLower)
		
		// Calculate keyword density
		analysis.PrimaryKeywordDensity = (float64(analysis.PrimaryKeywordCount) / float64(wordCount)) * 100
		
		// Optimal keyword density: 1-2%
		if analysis.PrimaryKeywordDensity >= 1.0 && analysis.PrimaryKeywordDensity <= 2.0 {
			analysis.PrimaryKeywordScore = 100
			analysis.PrimaryKeywordStatus = "optimal"
		} else if analysis.PrimaryKeywordDensity >= 0.5 && analysis.PrimaryKeywordDensity < 3.0 {
			analysis.PrimaryKeywordScore = 80
			analysis.PrimaryKeywordStatus = "good"
		} else if analysis.PrimaryKeywordDensity < 0.5 {
			analysis.PrimaryKeywordScore = 50
			analysis.PrimaryKeywordStatus = "too_low"
		} else {
			analysis.PrimaryKeywordScore = 30
			analysis.PrimaryKeywordStatus = "too_high"
		}

		// Check primary keyword in first 100 words
		firstHundredWords := sa.getFirstNWords(contentText, 100)
		if strings.Contains(strings.ToLower(firstHundredWords), primaryKeywordLower) {
			analysis.PrimaryKeywordInIntro = true
			analysis.IntroKeywordScore = 100
		} else {
			analysis.IntroKeywordScore = 40
		}
	}

	// Analyze secondary keywords
	for _, secondaryKeyword := range content.SecondaryKeywords {
		if secondaryKeyword == "" {
			continue
		}
		
		secondaryKeywordLower := strings.ToLower(secondaryKeyword)
		count := strings.Count(contentLower, secondaryKeywordLower)
		density := (float64(count) / float64(wordCount)) * 100
		
		analysis.SecondaryKeywordData = append(analysis.SecondaryKeywordData, SecondaryKeywordData{
			Keyword: secondaryKeyword,
			Count:   count,
			Density: density,
		})
	}

	// Calculate LSI (Latent Semantic Indexing) keyword score
	analysis.LSIKeywords = sa.identifyLSIKeywords(contentText, content.PrimaryKeyword)
	analysis.LSIScore = math.Min(float64(len(analysis.LSIKeywords))*10, 100)

	return analysis
}

// analyzeReadability analyzes content readability
func (sa *SEOAnalyzer) analyzeReadability(content string) ReadabilityAnalysis {
	analysis := ReadabilityAnalysis{}

	sentences := sa.countSentences(content)
	words := sa.countWords(content)
	syllables := sa.countSyllables(content)

	if sentences == 0 || words == 0 {
		return analysis
	}

	analysis.SentenceCount = sentences
	analysis.WordCount = words
	analysis.SyllableCount = syllables
	analysis.AvgWordsPerSentence = float64(words) / float64(sentences)
	analysis.AvgSyllablesPerWord = float64(syllables) / float64(words)

	// Calculate Flesch Reading Ease Score
	fleschScore := 206.835 - (1.015 * analysis.AvgWordsPerSentence) - (84.6 * analysis.AvgSyllablesPerWord)
	analysis.FleschScore = math.Max(0, math.Min(100, fleschScore))

	// Determine reading level
	if analysis.FleschScore >= 90 {
		analysis.ReadingLevel = "very_easy"
		analysis.ReadabilityScore = 100
	} else if analysis.FleschScore >= 80 {
		analysis.ReadingLevel = "easy"
		analysis.ReadabilityScore = 90
	} else if analysis.FleschScore >= 70 {
		analysis.ReadingLevel = "fairly_easy"
		analysis.ReadabilityScore = 80
	} else if analysis.FleschScore >= 60 {
		analysis.ReadingLevel = "standard"
		analysis.ReadabilityScore = 70
	} else if analysis.FleschScore >= 50 {
		analysis.ReadingLevel = "fairly_difficult"
		analysis.ReadabilityScore = 60
	} else if analysis.FleschScore >= 30 {
		analysis.ReadingLevel = "difficult"
		analysis.ReadabilityScore = 40
	} else {
		analysis.ReadingLevel = "very_difficult"
		analysis.ReadabilityScore = 20
	}

	// Calculate Flesch-Kincaid Grade Level
	analysis.FleschKincaidGrade = (0.39 * analysis.AvgWordsPerSentence) + (11.8 * analysis.AvgSyllablesPerWord) - 15.59

	// Analyze sentence length distribution
	analysis.SentenceLengthAnalysis = sa.analyzeSentenceLengths(content)

	// Check for transition words
	transitionWords := []string{"however", "therefore", "furthermore", "moreover", "additionally", "consequently", "meanwhile", "nevertheless", "similarly", "in contrast", "on the other hand", "in addition", "for example", "for instance"}
	contentLower := strings.ToLower(content)
	
	for _, word := range transitionWords {
		if strings.Contains(contentLower, word) {
			analysis.TransitionWords = append(analysis.TransitionWords, word)
		}
	}
	
	analysis.TransitionWordScore = math.Min(float64(len(analysis.TransitionWords))*15, 100)

	return analysis
}

// analyzeTechnicalSEO analyzes technical SEO factors
func (sa *SEOAnalyzer) analyzeTechnicalSEO(content ContentData) TechnicalAnalysis {
	analysis := TechnicalAnalysis{}

	// Analyze URL structure
	analysis.URLAnalysis = sa.analyzeURL(content.URL, content.PrimaryKeyword)

	// Analyze schema markup
	analysis.HasSchemaMarkup = content.SchemaMarkup != ""
	if analysis.HasSchemaMarkup {
		analysis.SchemaScore = 100
	} else {
		analysis.SchemaScore = 0
	}

	// Analyze canonical URL
	analysis.HasCanonicalURL = content.CanonicalURL != ""
	if analysis.HasCanonicalURL {
		analysis.CanonicalScore = 100
	} else {
		analysis.CanonicalScore = 50
	}

	// Analyze load time (if available)
	if content.LoadTime > 0 {
		if content.LoadTime <= 2.0 {
			analysis.LoadTimeScore = 100
			analysis.LoadTimeStatus = "excellent"
		} else if content.LoadTime <= 3.0 {
			analysis.LoadTimeScore = 80
			analysis.LoadTimeStatus = "good"
		} else if content.LoadTime <= 5.0 {
			analysis.LoadTimeScore = 60
			analysis.LoadTimeStatus = "fair"
		} else {
			analysis.LoadTimeScore = 30
			analysis.LoadTimeStatus = "poor"
		}
	} else {
		analysis.LoadTimeScore = 50 // Neutral score if no data
	}

	// Analyze mobile responsiveness
	analysis.IsMobileResponsive = content.MobileResponsive
	if analysis.IsMobileResponsive {
		analysis.MobileScore = 100
	} else {
		analysis.MobileScore = 0
	}

	return analysis
}

// analyzeLinkStructure analyzes internal and external link structure
func (sa *SEOAnalyzer) analyzeLinkStructure(internalLinks, externalLinks []LinkData) LinkAnalysis {
	analysis := LinkAnalysis{
		InternalLinks: internalLinks,
		ExternalLinks: externalLinks,
	}

	analysis.InternalLinkCount = len(internalLinks)
	analysis.ExternalLinkCount = len(externalLinks)

	// Analyze internal link optimization (3-10 internal links is optimal)
	if analysis.InternalLinkCount >= 3 && analysis.InternalLinkCount <= 10 {
		analysis.InternalLinkScore = 100
		analysis.InternalLinkStatus = "optimal"
	} else if analysis.InternalLinkCount >= 1 && analysis.InternalLinkCount <= 15 {
		analysis.InternalLinkScore = 80
		analysis.InternalLinkStatus = "good"
	} else if analysis.InternalLinkCount == 0 {
		analysis.InternalLinkScore = 20
		analysis.InternalLinkStatus = "missing"
	} else {
		analysis.InternalLinkScore = 60
		analysis.InternalLinkStatus = "too_many"
	}

	// Analyze external link optimization (2-5 external links is good)
	if analysis.ExternalLinkCount >= 2 && analysis.ExternalLinkCount <= 5 {
		analysis.ExternalLinkScore = 100
		analysis.ExternalLinkStatus = "optimal"
	} else if analysis.ExternalLinkCount >= 1 && analysis.ExternalLinkCount <= 8 {
		analysis.ExternalLinkScore = 80
		analysis.ExternalLinkStatus = "good"
	} else if analysis.ExternalLinkCount == 0 {
		analysis.ExternalLinkScore = 40
		analysis.ExternalLinkStatus = "none"
	} else {
		analysis.ExternalLinkScore = 60
		analysis.ExternalLinkStatus = "too_many"
	}

	// Check for proper anchor text usage
	analysis.AnchorTextAnalysis = sa.analyzeAnchorTexts(append(internalLinks, externalLinks...))

	return analysis
}

// analyzeImageOptimization analyzes image SEO optimization
func (sa *SEOAnalyzer) analyzeImageOptimization(images []ImageData) ImageAnalysis {
	analysis := ImageAnalysis{
		Images:     images,
		ImageCount: len(images),
	}

	if len(images) == 0 {
		analysis.ImageScore = 50 // Neutral score for content without images
		return analysis
	}

	altTextCount := 0
	titleCount := 0
	optimizedFileNames := 0

	for _, image := range images {
		if image.AltText != "" {
			altTextCount++
		}
		if image.Title != "" {
			titleCount++
		}
		if sa.isOptimizedFileName(image.FileName) {
			optimizedFileNames++
		}
	}

	analysis.ImagesWithAltText = altTextCount
	analysis.ImagesWithTitle = titleCount
	analysis.OptimizedFileNames = optimizedFileNames

	// Calculate scores
	altTextPercentage := float64(altTextCount) / float64(len(images)) * 100
	titlePercentage := float64(titleCount) / float64(len(images)) * 100
	fileNamePercentage := float64(optimizedFileNames) / float64(len(images)) * 100

	analysis.AltTextScore = altTextPercentage
	analysis.TitleScore = titlePercentage
	analysis.FileNameScore = fileNamePercentage

	// Overall image optimization score
	analysis.ImageScore = (analysis.AltTextScore*0.5 + analysis.TitleScore*0.2 + analysis.FileNameScore*0.3)

	return analysis
}

// calculateOverallSEOScore calculates the overall SEO score
func (sa *SEOAnalyzer) calculateOverallSEOScore(analysis SEOAnalysis) int {
	var totalScore float64
	var weights float64

	// Title optimization (15% weight)
	titleScore := (float64(analysis.TitleAnalysis.LengthScore) + float64(analysis.TitleAnalysis.KeywordScore)) / 2
	totalScore += titleScore * 0.15
	weights += 0.15

	// Meta description (10% weight)
	metaScore := (float64(analysis.MetaAnalysis.LengthScore) + float64(analysis.MetaAnalysis.KeywordScore)) / 2
	totalScore += metaScore * 0.10
	weights += 0.10

	// Content structure (20% weight)
	structureScore := (float64(analysis.StructureAnalysis.WordCountScore) + float64(analysis.StructureAnalysis.H1Score) + float64(analysis.StructureAnalysis.H2Score)) / 3
	totalScore += structureScore * 0.20
	weights += 0.20

	// Keyword optimization (20% weight)
	keywordScore := float64(analysis.KeywordAnalysis.PrimaryKeywordScore)
	totalScore += keywordScore * 0.20
	weights += 0.20

	// Readability (15% weight)
	readabilityScore := float64(analysis.ReadabilityAnalysis.ReadabilityScore)
	totalScore += readabilityScore * 0.15
	weights += 0.15

	// Technical SEO (10% weight)
	technicalScore := (float64(analysis.TechnicalAnalysis.SchemaScore) + float64(analysis.TechnicalAnalysis.CanonicalScore) + float64(analysis.TechnicalAnalysis.LoadTimeScore)) / 3
	totalScore += technicalScore * 0.10
	weights += 0.10

	// Link structure (5% weight)
	linkScore := (float64(analysis.LinkAnalysis.InternalLinkScore) + float64(analysis.LinkAnalysis.ExternalLinkScore)) / 2
	totalScore += linkScore * 0.05
	weights += 0.05

	// Image optimization (5% weight)
	imageScore := analysis.ImageAnalysis.ImageScore
	totalScore += imageScore * 0.05
	weights += 0.05

	if weights > 0 {
		return int(math.Min(totalScore/weights, 100))
	}

	return 0
}

// Helper methods

func (sa *SEOAnalyzer) countWords(text string) int {
	words := regexp.MustCompile(`\S+`).FindAllString(text, -1)
	return len(words)
}

func (sa *SEOAnalyzer) countSentences(text string) int {
	sentences := regexp.MustCompile(`[.!?]+`).Split(text, -1)
	return len(sentences) - 1 // Last split is usually empty
}

func (sa *SEOAnalyzer) countSyllables(text string) int {
	// Simplified syllable counting
	words := strings.Fields(text)
	totalSyllables := 0
	
	for _, word := range words {
		word = strings.ToLower(regexp.MustCompile(`[^a-z]`).ReplaceAllString(word, ""))
		if word == "" {
			continue
		}
		
		syllables := 0
		vowels := "aeiouy"
		prevWasVowel := false
		
		for _, char := range word {
			isVowel := strings.ContainsRune(vowels, char)
			if isVowel && !prevWasVowel {
				syllables++
			}
			prevWasVowel = isVowel
		}
		
		// Silent 'e' rule
		if strings.HasSuffix(word, "e") && syllables > 1 {
			syllables--
		}
		
		// Minimum 1 syllable per word
		if syllables == 0 {
			syllables = 1
		}
		
		totalSyllables += syllables
	}
	
	return totalSyllables
}

func (sa *SEOAnalyzer) getFirstNWords(text string, n int) string {
	words := strings.Fields(text)
	if len(words) <= n {
		return text
	}
	return strings.Join(words[:n], " ")
}

func (sa *SEOAnalyzer) identifyLSIKeywords(content, primaryKeyword string) []string {
	// Simplified LSI keyword identification
	// In a real implementation, this would use more sophisticated NLP
	primaryLower := strings.ToLower(primaryKeyword)
	contentLower := strings.ToLower(content)
	
	// Define some common LSI patterns based on primary keyword
	lsiKeywords := []string{}
	
	// This is a simplified example - real LSI would be much more sophisticated
	if strings.Contains(primaryLower, "seo") {
		potentialLSI := []string{"search engine", "optimization", "ranking", "keywords", "google", "content marketing", "backlinks", "meta tags"}
		for _, lsi := range potentialLSI {
			if strings.Contains(contentLower, lsi) {
				lsiKeywords = append(lsiKeywords, lsi)
			}
		}
	}
	
	return lsiKeywords
}

func (sa *SEOAnalyzer) analyzeSentenceLengths(content string) SentenceLengthAnalysis {
	sentences := regexp.MustCompile(`[.!?]+`).Split(content, -1)
	analysis := SentenceLengthAnalysis{}
	
	var lengths []int
	for _, sentence := range sentences {
		words := strings.Fields(strings.TrimSpace(sentence))
		if len(words) > 0 {
			lengths = append(lengths, len(words))
		}
	}
	
	if len(lengths) == 0 {
		return analysis
	}
	
	// Calculate statistics
	sum := 0
	for _, length := range lengths {
		sum += length
		if length > analysis.LongestSentence {
			analysis.LongestSentence = length
		}
		if analysis.ShortestSentence == 0 || length < analysis.ShortestSentence {
			analysis.ShortestSentence = length
		}
	}
	
	analysis.AverageLength = float64(sum) / float64(len(lengths))
	
	// Count sentences by length category
	for _, length := range lengths {
		if length <= 10 {
			analysis.ShortSentences++
		} else if length <= 20 {
			analysis.MediumSentences++
		} else {
			analysis.LongSentences++
		}
	}
	
	return analysis
}

func (sa *SEOAnalyzer) analyzeURL(url, primaryKeyword string) URLAnalysis {
	analysis := URLAnalysis{
		URL:    url,
		Length: len(url),
	}
	
	// Check URL length (optimal: under 75 characters)
	if analysis.Length <= 75 {
		analysis.LengthScore = 100
		analysis.LengthStatus = "optimal"
	} else if analysis.Length <= 100 {
		analysis.LengthScore = 80
		analysis.LengthStatus = "good"
	} else {
		analysis.LengthScore = 60
		analysis.LengthStatus = "too_long"
	}
	
	// Check for keyword in URL
	if primaryKeyword != "" {
		urlLower := strings.ToLower(url)
		keywordLower := strings.ToLower(primaryKeyword)
		keywordSlug := strings.ReplaceAll(keywordLower, " ", "-")
		
		if strings.Contains(urlLower, keywordSlug) || strings.Contains(urlLower, keywordLower) {
			analysis.ContainsKeyword = true
			analysis.KeywordScore = 100
		} else {
			analysis.KeywordScore = 30
		}
	}
	
	// Check URL structure
	if regexp.MustCompile(`^https?://[^/]+/[a-z0-9-]+/?$`).MatchString(strings.ToLower(url)) {
		analysis.StructureScore = 100
		analysis.Structure = "clean"
	} else {
		analysis.StructureScore = 70
		analysis.Structure = "complex"
	}
	
	return analysis
}

func (sa *SEOAnalyzer) analyzeAnchorTexts(links []LinkData) AnchorTextAnalysis {
	analysis := AnchorTextAnalysis{}
	
	if len(links) == 0 {
		return analysis
	}
	
	anchorTexts := make(map[string]int)
	totalLinks := len(links)
	
	for _, link := range links {
		anchorText := strings.ToLower(strings.TrimSpace(link.AnchorText))
		if anchorText != "" {
			anchorTexts[anchorText]++
		}
	}
	
	// Check for over-optimization (same anchor text used too frequently)
	maxFrequency := 0
	for anchorText, count := range anchorTexts {
		frequency := (count * 100) / totalLinks
		if frequency > maxFrequency {
			maxFrequency = frequency
			analysis.MostUsedAnchorText = anchorText
		}
		
		if frequency > 30 { // More than 30% is over-optimization
			analysis.OverOptimizedAnchors = append(analysis.OverOptimizedAnchors, anchorText)
		}
	}
	
	analysis.AnchorTextVariety = len(anchorTexts)
	analysis.MaxAnchorFrequency = maxFrequency
	
	// Calculate diversity score
	if totalLinks > 0 {
		analysis.DiversityScore = math.Min(float64(analysis.AnchorTextVariety)/float64(totalLinks)*100, 100)
	}
	
	return analysis
}

func (sa *SEOAnalyzer) isOptimizedFileName(fileName string) bool {
	// Check if filename contains descriptive words and uses hyphens
	fileName = strings.ToLower(fileName)
	
	// Remove file extension
	if dotIndex := strings.LastIndex(fileName, "."); dotIndex > 0 {
		fileName = fileName[:dotIndex]
	}
	
	// Check for descriptive content (not just numbers or generic names)
	if regexp.MustCompile(`^(img|image|picture|photo)\d*$`).MatchString(fileName) {
		return false
	}
	
	// Check for proper formatting (uses hyphens, no underscores or spaces)
	if strings.Contains(fileName, "_") || strings.Contains(fileName, " ") {
		return false
	}
	
	// Check for meaningful content (at least 3 characters and contains letters)
	if len(fileName) < 3 || !regexp.MustCompile(`[a-z]`).MatchString(fileName) {
		return false
	}
	
	return true
}

// generateRecommendations generates actionable SEO recommendations
func (sa *SEOAnalyzer) generateRecommendations(analysis SEOAnalysis) []string {
	var recommendations []string
	
	// Title recommendations
	if analysis.TitleAnalysis.LengthStatus == "too_short" {
		recommendations = append(recommendations, "Expand your title to 50-60 characters for optimal search engine display")
	} else if analysis.TitleAnalysis.LengthStatus == "too_long" {
		recommendations = append(recommendations, "Shorten your title to under 60 characters to avoid truncation in search results")
	}
	
	if !analysis.TitleAnalysis.ContainsPrimaryKeyword {
		recommendations = append(recommendations, "Include your primary keyword in the title, preferably near the beginning")
	}
	
	// Meta description recommendations
	if analysis.MetaAnalysis.LengthStatus == "too_short" {
		recommendations = append(recommendations, "Expand your meta description to 150-160 characters for better search result display")
	} else if analysis.MetaAnalysis.LengthStatus == "too_long" {
		recommendations = append(recommendations, "Shorten your meta description to under 160 characters")
	}
	
	if !analysis.MetaAnalysis.CallToAction {
		recommendations = append(recommendations, "Add a compelling call-to-action to your meta description")
	}
	
	// Content structure recommendations
	if analysis.StructureAnalysis.H1Status == "missing" {
		recommendations = append(recommendations, "Add an H1 heading to your content for better structure")
	} else if analysis.StructureAnalysis.H1Status == "multiple" {
		recommendations = append(recommendations, "Use only one H1 heading per page")
	}
	
	if analysis.StructureAnalysis.H2Status == "missing" {
		recommendations = append(recommendations, "Add H2 subheadings to improve content structure and readability")
	}
	
	// Keyword recommendations
	if analysis.KeywordAnalysis.PrimaryKeywordStatus == "too_low" {
		recommendations = append(recommendations, "Increase primary keyword usage to 1-2% density")
	} else if analysis.KeywordAnalysis.PrimaryKeywordStatus == "too_high" {
		recommendations = append(recommendations, "Reduce primary keyword usage to avoid over-optimization (aim for 1-2% density)")
	}
	
	// Readability recommendations
	if analysis.ReadabilityAnalysis.FleschScore < 60 {
		recommendations = append(recommendations, "Improve readability by using shorter sentences and simpler words")
	}
	
	// Technical recommendations
	if !analysis.TechnicalAnalysis.HasSchemaMarkup {
		recommendations = append(recommendations, "Add schema markup to help search engines understand your content better")
	}
	
	// Link recommendations
	if analysis.LinkAnalysis.InternalLinkStatus == "missing" {
		recommendations = append(recommendations, "Add 3-5 internal links to related content on your website")
	}
	
	if analysis.LinkAnalysis.ExternalLinkStatus == "none" {
		recommendations = append(recommendations, "Include 2-3 links to high-quality external sources for credibility")
	}
	
	// Image recommendations
	if analysis.ImageAnalysis.AltTextScore < 80 {
		recommendations = append(recommendations, "Add descriptive alt text to all images for better accessibility and SEO")
	}
	
	return recommendations
}

// identifyOpportunities identifies specific optimization opportunities
func (sa *SEOAnalyzer) identifyOpportunities(analysis SEOAnalysis) []Opportunity {
	var opportunities []Opportunity
	
	// High-impact opportunities
	if !analysis.TitleAnalysis.ContainsPrimaryKeyword {
		opportunities = append(opportunities, Opportunity{
			Type:        "title_optimization",
			Priority:    "high",
			Impact:      "high",
			Effort:      "low",
			Title:       "Add primary keyword to title",
			Description: "Including your primary keyword in the title can significantly improve rankings",
			Action:      "Edit the title to include your primary keyword, preferably at the beginning",
		})
	}
	
	if analysis.KeywordAnalysis.PrimaryKeywordDensity < 0.5 {
		opportunities = append(opportunities, Opportunity{
			Type:        "keyword_optimization",
			Priority:    "high",
			Impact:      "medium",
			Effort:      "medium",
			Title:       "Increase keyword density",
			Description: "Your primary keyword density is too low, which may affect rankings",
			Action:      "Naturally incorporate your primary keyword 2-3 more times in the content",
		})
	}
	
	// Medium-impact opportunities
	if analysis.StructureAnalysis.H2Count < 2 {
		opportunities = append(opportunities, Opportunity{
			Type:        "content_structure",
			Priority:    "medium",
			Impact:      "medium",
			Effort:      "low",
			Title:       "Add more H2 subheadings",
			Description: "Proper heading structure improves both SEO and user experience",
			Action:      "Break your content into logical sections with descriptive H2 headings",
		})
	}
	
	if len(analysis.LinkAnalysis.InternalLinks) < 3 {
		opportunities = append(opportunities, Opportunity{
			Type:        "internal_linking",
			Priority:    "medium",
			Impact:      "medium",
			Effort:      "low",
			Title:       "Add internal links",
			Description: "Internal links help search engines understand your site structure and keep users engaged",
			Action:      "Link to 3-5 related articles or pages on your website",
		})
	}
	
	// Low-impact but easy wins
	if analysis.ImageAnalysis.AltTextScore < 100 {
		opportunities = append(opportunities, Opportunity{
			Type:        "image_optimization",
			Priority:    "low",
			Impact:      "low",
			Effort:      "low",
			Title:       "Optimize image alt text",
			Description: "Alt text helps search engines understand your images and improves accessibility",
			Action:      "Add descriptive alt text to all images, including relevant keywords where appropriate",
		})
	}
	
	if !analysis.TechnicalAnalysis.HasSchemaMarkup {
		opportunities = append(opportunities, Opportunity{
			Type:        "schema_markup",
			Priority:    "low",
			Impact:      "medium",
			Effort:      "high",
			Title:       "Add schema markup",
			Description: "Schema markup helps search engines better understand and display your content",
			Action:      "Implement appropriate schema markup (Article, BlogPosting, etc.) for your content type",
		})
	}
	
	return opportunities
}