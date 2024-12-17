-- Sample data for full-text search benchmarking

-- Create documents table with full-text search optimization
CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample technical documents
INSERT INTO documents (title, content) VALUES 
(
    'Advances in Machine Learning', 
    'Machine learning continues to revolutionize artificial intelligence through deep neural networks, reinforcement learning, and advanced algorithmic approaches. Modern ML techniques are transforming industries from healthcare to finance.'
),
(
    'Cybersecurity in the Digital Age', 
    'As technology evolves, cybersecurity becomes increasingly critical. Emerging threats require sophisticated defense mechanisms, including AI-powered threat detection, zero-trust architectures, and advanced encryption techniques.'
),
(
    'Cloud Computing Trends', 
    'Cloud computing is rapidly changing enterprise infrastructure. Serverless architectures, containerization with Kubernetes, and multi-cloud strategies are becoming mainstream approaches for scalable and flexible IT solutions.'
),
(
    'Data Science and Big Data Analytics', 
    'Big data analytics provides unprecedented insights across industries. Techniques like predictive modeling, machine learning, and advanced statistical analysis help organizations make data-driven decisions.'
),
(
    'Blockchain Technology Overview', 
    'Blockchain technology extends beyond cryptocurrencies, offering decentralized solutions for secure, transparent transactions. Smart contracts and distributed ledger technologies are reshaping finance, supply chain, and governance models.'
);

-- Create full-text search index
CREATE INDEX IF NOT EXISTS idx_documents_search 
ON documents 
USING GIN (to_tsvector('english', title || ' ' || content));