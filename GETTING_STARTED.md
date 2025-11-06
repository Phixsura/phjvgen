# Getting Started with phjvgen

This quick guide will help you get started with phjvgen in 5 minutes.

## Prerequisites

- Go 1.24+ (for building)
- Java 25 (for running generated projects)
- Maven 3.x (for building generated projects)
- MySQL 8.0+ (optional, for demo with database)

## Quick Start

### Step 1: Build phjvgen

```bash
# Navigate to phjvgen directory
cd phjvgen

# Build the binary
go build -o phjvgen .

# Verify it works
./phjvgen version
```

Expected output:
```
phjvgen version 1.0.0
Java 25 LTS Project Generator
Copyright (c) 2025
```

### Step 2: Generate Your First Project

#### Option A: Quick Example (No Input Required)

```bash
# Generate a pre-configured demo project
./phjvgen example

# Result: demo-app/ directory created with:
# - Group ID: com.example.demo
# - Artifact ID: demo-app
# - Complete project structure
```

#### Option B: Custom Project (Interactive)

```bash
# Start interactive project generation
./phjvgen generate

# You'll be prompted for:
# 1. Group ID (e.g., com.mycompany)
# 2. Artifact ID (e.g., my-awesome-app)
# 3. Version (default: 1.0.0)
# 4. Project Name (default: artifact ID)
# 5. Description
# 6. Output Directory (default: ./artifact-id)
```

### Step 3: Explore the Generated Project

```bash
# Enter the project
cd demo-app  # or your custom project name

# View structure
ls -la
```

You'll see:
```
demo-app/
├── pom.xml                    # Parent POM
├── common/                    # Common utilities
├── domain/                    # Domain entities
├── infrastructure/            # Data access
├── adapter/
│   ├── adapter-rest/         # REST controllers
│   └── adapter-schedule/     # Scheduled jobs
├── application/
│   └── application-user/     # User business logic
├── starter/                   # Application startup
├── README.md
└── .gitignore
```

### Step 4: Build the Project

```bash
# Build with Maven
mvn clean install
```

This compiles all modules and creates the executable JAR.

### Step 5: Run the Application

```bash
# Run the application
java --enable-preview -jar starter/target/starter-1.0.0.jar
```

In another terminal:

```bash
# Test the health endpoint
curl http://localhost:8080/api/health

# Expected response:
# {"code":200,"message":"success","data":"Service is running","timestamp":...}
```

## Next Steps

### Add CRUD Demo Code (Optional)

If you want a complete working example with database:

```bash
# Generate User CRUD example
../phjvgen demo

# This creates:
# - User entity and repository
# - User service and DTOs
# - User REST controller
# - Database migration SQL
# - Exception handlers
```

### Add More Business Modules

```bash
# Add a payment module
../phjvgen add payment

# Add an order module
../phjvgen add order

# Each module gets:
# - Complete directory structure
# - Module POM
# - Sample service class
# - Auto-updated parent POM
```

### Setup Database (For Demo Code)

If you generated demo code:

```bash
# Create database
mysql -u root -p
> CREATE DATABASE demo_app;
> USE demo_app;
> source infrastructure/src/main/resources/db/migration/V1__create_user_table.sql;
> exit;

# Configure database connection
# Edit: starter/src/main/resources/application-dev.yml

# Restart application
mvn clean install
java --enable-preview -jar starter/target/starter-1.0.0.jar
```

Test the User API:

```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"username":"john","email":"john@example.com","phone":"1234567890"}'

# Get all users
curl http://localhost:8080/api/users

# Get user by ID
curl http://localhost:8080/api/users/1

# Update user
curl -X PUT http://localhost:8080/api/users/1 \
  -H 'Content-Type: application/json' \
  -d '{"email":"newemail@example.com","status":1}'

# Delete user
curl -X DELETE http://localhost:8080/api/users/1
```

## Installation (Optional)

Install phjvgen system-wide:

```bash
# Install to system
./phjvgen install

# Now use from anywhere
phjvgen version
phjvgen generate
```

Installation locations:
- **Linux/macOS (user)**: `~/.local/bin/phjvgen`
- **Linux/macOS (root)**: `/usr/local/bin/phjvgen`

Add to PATH if needed:
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

## Common Commands

```bash
# View all commands
phjvgen --help

# View command-specific help
phjvgen generate --help
phjvgen demo --help
phjvgen add --help

# Generate new project
phjvgen generate    # or: phjvgen gen, phjvgen g

# Generate quick example
phjvgen example

# Add CRUD demo (run in project directory)
phjvgen demo

# Add new module (run in project directory)
phjvgen add <module-name>

# Install to system
phjvgen install

# Check version
phjvgen version
```

## Project Structure Overview

```
your-project/
├── common/              # Shared code (exceptions, responses, constants)
├── domain/              # Business entities and repository interfaces
├── infrastructure/      # Database access, caching, external services
├── adapter/
│   ├── adapter-rest/   # HTTP controllers, request/response DTOs
│   └── adapter-schedule/ # Scheduled tasks
├── application/         # Business logic, use cases
│   └── application-user/ # Example: User business module
└── starter/            # Spring Boot application entry point
```

**Layer Dependencies** (DDD):
```
Starter → Adapter → Application → Domain ← Infrastructure
                                     ↓
                                  Common
```

## Tips

1. **Start Simple**: Use `phjvgen example` first to see the structure
2. **Add Demo Code**: Run `phjvgen demo` to see a complete CRUD example
3. **Customize**: Modify generated code to fit your needs
4. **Add Modules**: Use `phjvgen add` for new business domains
5. **Follow DDD**: Keep domain logic in domain layer, infrastructure details in infrastructure layer

## Troubleshooting

### Build fails with "go: command not found"

Install Go from https://golang.org/dl/

### Generated project build fails

Check Java and Maven versions:
```bash
java --version    # Should be 25+
mvn --version     # Should be 3.6+
```

### Application fails to start

Check if port 8080 is available:
```bash
lsof -i :8080
# Kill process if needed: kill <PID>
```

### Database connection error (demo code)

1. Check MySQL is running
2. Verify database exists
3. Update credentials in `starter/src/main/resources/application-dev.yml`

## Learn More

- **Full Documentation**: See [README.md](README.md)
- **Build Guide**: See [BUILD_AND_TEST.md](BUILD_AND_TEST.md)
- **Refactoring Details**: See [../REFACTORING_SUMMARY.md](../REFACTORING_SUMMARY.md)

## Need Help?

1. Run `phjvgen --help`
2. Check command-specific help: `phjvgen <command> --help`
3. Review generated README.md in your project
4. Read the documentation files

## Example Session

Complete example session from scratch:

```bash
# 1. Build phjvgen
cd phjvgen
go build -o phjvgen .

# 2. Generate example project
./phjvgen example

# 3. Add demo code
cd demo-app
../phjvgen demo
# Answer 'y' to confirm

# 4. Add payment module
../phjvgen add payment
# Answer 'y' to confirm

# 5. Build project
mvn clean install

# 6. Run application
java --enable-preview -jar starter/target/starter-1.0.0.jar &

# 7. Test
curl http://localhost:8080/api/health

# 8. Stop application
pkill -f starter-1.0.0.jar
```

Congratulations! You're now ready to use phjvgen to generate Java projects.
