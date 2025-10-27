# Python 2.7 Dockerfile (Deprecated - for demonstration purposes only)
FROM python:2.7-slim

WORKDIR /app

# Fix Debian Buster EOL - use archived repositories
RUN sed -i 's|http://deb.debian.org|http://archive.debian.org|g' /etc/apt/sources.list && \
    sed -i 's|http://security.debian.org|http://archive.debian.org|g' /etc/apt/sources.list && \
    sed -i '/stretch-updates/d' /etc/apt/sources.list

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Make scripts executable
RUN chmod +x app.py

# Expose REST port
EXPOSE 8083

# Run the REST application
CMD ["python", "app.py"]
