package tests

// SampleArticleHTML provides a realistic HTML page for testing content extraction
const SampleArticleHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Test Article Page</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="author" content="John Doe">
    <meta name="description" content="This is a test article for content extraction">
</head>
<body>
    <header>
        <nav>
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="/articles">Articles</a></li>
                <li><a href="/about">About</a></li>
                <li><a href="/contact">Contact</a></li>
            </ul>
        </nav>
        <div class="site-branding">
            <h1 class="site-title">Test News Site</h1>
            <p class="site-description">Your source for test articles</p>
        </div>
    </header>

    <main>
        <article class="post">
            <header class="entry-header">
                <h1 class="entry-title">Understanding Content Extraction</h1>
                <div class="entry-meta">
                    <span class="posted-on">Posted on <time datetime="2023-04-15">April 15, 2023</time></span>
                    <span class="byline">by <span class="author">Jane Smith</span></span>
                </div>
            </header>

            <div class="entry-content">
                <p>Content extraction is a critical part of web scraping and data processing. It allows us to identify and extract the main content from a webpage, filtering out all the surrounding navigation, advertisements, and other irrelevant elements.</p>
                
                <h2>Why Content Extraction Matters</h2>
                <p>When building a content aggregator or reader application, users want to see the actual content of articles without all the clutter. Content extraction algorithms help achieve this goal by:
                </p>
                <ul>
                    <li>Identifying the main article content</li>
                    <li>Removing navigation elements</li>
                    <li>Eliminating advertisements</li>
                    <li>Preserving important formatting and media</li>
                </ul>
                
                <p>Many content extraction libraries use heuristics to identify the main content area of a webpage. These heuristics can include:</p>
                
                <h3>Density-Based Approaches</h3>
                <p>These approaches look at the ratio of text to HTML tags. Areas with a high density of text are likely to be the main content.</p>
                
                <h3>DOM-Based Approaches</h3>
                <p>These approaches analyze the Document Object Model (DOM) structure to find patterns that typically indicate content areas.</p>
                
                <img src="/images/content-extraction.jpg" alt="Content Extraction Diagram" />
                
                <p>The table below compares some popular content extraction libraries:</p>
                
                <table>
                    <thead>
                        <tr>
                            <th>Library</th>
                            <th>Language</th>
                            <th>Approach</th>
                            <th>Accuracy</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>Readability</td>
                            <td>JavaScript</td>
                            <td>Heuristic</td>
                            <td>High</td>
                        </tr>
                        <tr>
                            <td>Goose</td>
                            <td>Java/Go</td>
                            <td>ML + Heuristic</td>
                            <td>Medium-High</td>
                        </tr>
                        <tr>
                            <td>Newspaper</td>
                            <td>Python</td>
                            <td>ML + NLP</td>
                            <td>High</td>
                        </tr>
                    </tbody>
                </table>
                
                <h2>Conclusion</h2>
                <p>Effective content extraction is essential for any application that deals with web content. By using the right tools and approaches, we can provide users with a clean, distraction-free reading experience.</p>
            </div>
        </article>
    </main>

    <aside>
        <div class="widget">
            <h2>Related Articles</h2>
            <ul>
                <li><a href="/article1">Web Scraping Basics</a></li>
                <li><a href="/article2">HTML Parsing Techniques</a></li>
                <li><a href="/article3">Data Processing Pipelines</a></li>
            </ul>
        </div>
        
        <div class="widget">
            <h2>Advertisement</h2>
            <div class="ad-container">
                <a href="https://example.com/ad"><img src="/ads/banner.jpg" alt="Advertisement" /></a>
            </div>
        </div>
    </aside>

    <footer>
        <div class="footer-widgets">
            <div class="widget">
                <h3>About Us</h3>
                <p>Test News Site is a fictional website created for testing purposes.</p>
            </div>
            <div class="widget">
                <h3>Categories</h3>
                <ul>
                    <li><a href="/category/tech">Technology</a></li>
                    <li><a href="/category/science">Science</a></li>
                    <li><a href="/category/health">Health</a></li>
                </ul>
            </div>
        </div>
        <div class="site-info">
            <p>&copy; 2023 Test News Site. All rights reserved.</p>
        </div>
    </footer>

    <div id="cookie-notice">
        <p>This website uses cookies to improve your experience.</p>
        <button>Accept</button>
        <button>Decline</button>
    </div>
</body>
</html>`
