package templates

// ParentPOM is the parent POM template
const ParentPOM = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>{{GROUP_ID}}</groupId>
    <artifactId>{{ARTIFACT_ID}}</artifactId>
    <version>{{VERSION}}</version>
    <packaging>pom</packaging>

    <name>{{PROJECT_NAME}}</name>
    <description>{{PROJECT_DESCRIPTION}}</description>

    <modules>
        <module>common</module>
        <module>domain</module>
        <module>infrastructure</module>
        <module>adapter/adapter-rest</module>
        <module>adapter/adapter-schedule</module>
        <module>application/application-user</module>
        <module>starter</module>
    </modules>

    <properties>
        <!-- Java 25 LTS -->
        <java.version>25</java.version>
        <maven.compiler.source>25</maven.compiler.source>
        <maven.compiler.target>25</maven.compiler.target>
        <maven.compiler.release>25</maven.compiler.release>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
        <project.reporting.outputEncoding>UTF-8</project.reporting.outputEncoding>

        <!-- Spring Boot -->
        <spring-boot.version>4.0.0-RC1</spring-boot.version>

        <!-- 数据库 -->
        <mybatis-plus.version>3.5.8</mybatis-plus.version>
        <mysql.version>8.0.33</mysql.version>
        <hikaricp.version>6.0.0</hikaricp.version>

        <!-- 工具库 -->
        <lombok.version>1.18.42</lombok.version>
        <mapstruct.version>1.6.0</mapstruct.version>
        <hutool.version>5.8.28</hutool.version>
        <guava.version>33.3.0-jre</guava.version>
        <commons-lang3.version>3.15.0</commons-lang3.version>

        <!-- Redis -->
        <redisson.version>3.30.0</redisson.version>
        <caffeine.version>3.1.8</caffeine.version>

        <!-- JSON -->
        <jackson.version>2.17.0</jackson.version>
        <fastjson2.version>2.0.52</fastjson2.version>

        <!-- 监控 -->
        <micrometer.version>1.13.0</micrometer.version>

        <!-- 插件版本 -->
        <maven-compiler-plugin.version>3.13.0</maven-compiler-plugin.version>
        <maven-surefire-plugin.version>3.2.5</maven-surefire-plugin.version>
    </properties>

    <dependencyManagement>
        <dependencies>
            <!-- Spring Boot Dependencies -->
            <dependency>
                <groupId>org.springframework.boot</groupId>
                <artifactId>spring-boot-dependencies</artifactId>
                <version>${spring-boot.version}</version>
                <type>pom</type>
                <scope>import</scope>
            </dependency>

            <!-- 项目内部模块 -->
            <dependency>
                <groupId>{{GROUP_ID}}</groupId>
                <artifactId>common</artifactId>
                <version>${project.version}</version>
            </dependency>
            <dependency>
                <groupId>{{GROUP_ID}}</groupId>
                <artifactId>domain</artifactId>
                <version>${project.version}</version>
            </dependency>
            <dependency>
                <groupId>{{GROUP_ID}}</groupId>
                <artifactId>infrastructure</artifactId>
                <version>${project.version}</version>
            </dependency>
            <dependency>
                <groupId>{{GROUP_ID}}</groupId>
                <artifactId>adapter-rest</artifactId>
                <version>${project.version}</version>
            </dependency>
            <dependency>
                <groupId>{{GROUP_ID}}</groupId>
                <artifactId>adapter-schedule</artifactId>
                <version>${project.version}</version>
            </dependency>
            <dependency>
                <groupId>{{GROUP_ID}}</groupId>
                <artifactId>application-user</artifactId>
                <version>${project.version}</version>
            </dependency>

            <!-- MyBatis Plus -->
            <dependency>
                <groupId>com.baomidou</groupId>
                <artifactId>mybatis-plus-spring-boot3-starter</artifactId>
                <version>${mybatis-plus.version}</version>
            </dependency>

            <!-- MySQL -->
            <dependency>
                <groupId>com.mysql</groupId>
                <artifactId>mysql-connector-j</artifactId>
                <version>${mysql.version}</version>
            </dependency>

            <!-- HikariCP -->
            <dependency>
                <groupId>com.zaxxer</groupId>
                <artifactId>HikariCP</artifactId>
                <version>${hikaricp.version}</version>
            </dependency>

            <!-- Lombok -->
            <dependency>
                <groupId>org.projectlombok</groupId>
                <artifactId>lombok</artifactId>
                <version>${lombok.version}</version>
            </dependency>

            <!-- MapStruct -->
            <dependency>
                <groupId>org.mapstruct</groupId>
                <artifactId>mapstruct</artifactId>
                <version>${mapstruct.version}</version>
            </dependency>
            <dependency>
                <groupId>org.mapstruct</groupId>
                <artifactId>mapstruct-processor</artifactId>
                <version>${mapstruct.version}</version>
            </dependency>

            <!-- Hutool -->
            <dependency>
                <groupId>cn.hutool</groupId>
                <artifactId>hutool-all</artifactId>
                <version>${hutool.version}</version>
            </dependency>

            <!-- Guava -->
            <dependency>
                <groupId>com.google.guava</groupId>
                <artifactId>guava</artifactId>
                <version>${guava.version}</version>
            </dependency>

            <!-- Commons Lang3 -->
            <dependency>
                <groupId>org.apache.commons</groupId>
                <artifactId>commons-lang3</artifactId>
                <version>${commons-lang3.version}</version>
            </dependency>

            <!-- Redisson -->
            <dependency>
                <groupId>org.redisson</groupId>
                <artifactId>redisson-spring-boot-starter</artifactId>
                <version>${redisson.version}</version>
            </dependency>

            <!-- Caffeine -->
            <dependency>
                <groupId>com.github.ben-manes.caffeine</groupId>
                <artifactId>caffeine</artifactId>
                <version>${caffeine.version}</version>
            </dependency>

            <!-- FastJSON2 -->
            <dependency>
                <groupId>com.alibaba.fastjson2</groupId>
                <artifactId>fastjson2</artifactId>
                <version>${fastjson2.version}</version>
            </dependency>
        </dependencies>
    </dependencyManagement>

    <build>
        <pluginManagement>
            <plugins>
                <!-- Maven Compiler Plugin -->
                <plugin>
                    <groupId>org.apache.maven.plugins</groupId>
                    <artifactId>maven-compiler-plugin</artifactId>
                    <version>${maven-compiler-plugin.version}</version>
                    <configuration>
                        <release>${java.version}</release>
                        <compilerArgs>
                            <arg>--enable-preview</arg>
                        </compilerArgs>
                        <annotationProcessorPaths>
                            <path>
                                <groupId>org.projectlombok</groupId>
                                <artifactId>lombok</artifactId>
                                <version>${lombok.version}</version>
                            </path>
                            <path>
                                <groupId>org.mapstruct</groupId>
                                <artifactId>mapstruct-processor</artifactId>
                                <version>${mapstruct.version}</version>
                            </path>
                        </annotationProcessorPaths>
                    </configuration>
                </plugin>

                <!-- Maven Surefire Plugin -->
                <plugin>
                    <groupId>org.apache.maven.plugins</groupId>
                    <artifactId>maven-surefire-plugin</artifactId>
                    <version>${maven-surefire-plugin.version}</version>
                    <configuration>
                        <argLine>--enable-preview</argLine>
                    </configuration>
                </plugin>

                <!-- Spring Boot Maven Plugin -->
                <plugin>
                    <groupId>org.springframework.boot</groupId>
                    <artifactId>spring-boot-maven-plugin</artifactId>
                    <version>${spring-boot.version}</version>
                </plugin>
            </plugins>
        </pluginManagement>
    </build>

    <repositories>
        <repository>
            <id>central</id>
            <name>Maven Central</name>
            <url>https://repo.maven.apache.org/maven2</url>
        </repository>
    </repositories>
</project>
`

// CommonPOM is the common module POM template
const CommonPOM = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <parent>
        <groupId>{{GROUP_ID}}</groupId>
        <artifactId>{{ARTIFACT_ID}}</artifactId>
        <version>{{VERSION}}</version>
    </parent>

    <artifactId>common</artifactId>
    <packaging>jar</packaging>
    <name>common</name>
    <description>公共模块</description>

    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-validation</artifactId>
        </dependency>
        <dependency>
            <groupId>org.projectlombok</groupId>
            <artifactId>lombok</artifactId>
            <scope>provided</scope>
        </dependency>
        <dependency>
            <groupId>cn.hutool</groupId>
            <artifactId>hutool-all</artifactId>
        </dependency>
        <dependency>
            <groupId>com.fasterxml.jackson.core</groupId>
            <artifactId>jackson-databind</artifactId>
        </dependency>
    </dependencies>
</project>
`

// DomainPOM is the domain module POM template
const DomainPOM = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <parent>
        <groupId>{{GROUP_ID}}</groupId>
        <artifactId>{{ARTIFACT_ID}}</artifactId>
        <version>{{VERSION}}</version>
    </parent>

    <artifactId>domain</artifactId>
    <packaging>jar</packaging>
    <name>domain</name>
    <description>领域层</description>

    <dependencies>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>common</artifactId>
        </dependency>
        <dependency>
            <groupId>org.projectlombok</groupId>
            <artifactId>lombok</artifactId>
            <scope>provided</scope>
        </dependency>
    </dependencies>
</project>
`

// InfrastructurePOM is the infrastructure module POM template
const InfrastructurePOM = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <parent>
        <groupId>{{GROUP_ID}}</groupId>
        <artifactId>{{ARTIFACT_ID}}</artifactId>
        <version>{{VERSION}}</version>
    </parent>

    <artifactId>infrastructure</artifactId>
    <packaging>jar</packaging>
    <name>infrastructure</name>
    <description>基础设施层</description>

    <dependencies>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>domain</artifactId>
        </dependency>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>common</artifactId>
        </dependency>
        <dependency>
            <groupId>com.baomidou</groupId>
            <artifactId>mybatis-plus-spring-boot3-starter</artifactId>
        </dependency>
        <dependency>
            <groupId>com.mysql</groupId>
            <artifactId>mysql-connector-j</artifactId>
        </dependency>
        <dependency>
            <groupId>org.redisson</groupId>
            <artifactId>redisson-spring-boot-starter</artifactId>
            <optional>true</optional>
        </dependency>
        <dependency>
            <groupId>com.github.ben-manes.caffeine</groupId>
            <artifactId>caffeine</artifactId>
        </dependency>
        <dependency>
            <groupId>org.projectlombok</groupId>
            <artifactId>lombok</artifactId>
            <scope>provided</scope>
        </dependency>
    </dependencies>
</project>
`

// AdapterRestPOM is the adapter-rest module POM template
const AdapterRestPOM = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <parent>
        <groupId>{{GROUP_ID}}</groupId>
        <artifactId>{{ARTIFACT_ID}}</artifactId>
        <version>{{VERSION}}</version>
        <relativePath>../../pom.xml</relativePath>
    </parent>

    <artifactId>adapter-rest</artifactId>
    <packaging>jar</packaging>
    <name>adapter-rest</name>
    <description>REST适配器</description>

    <dependencies>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>application-user</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-validation</artifactId>
        </dependency>
        <dependency>
            <groupId>org.projectlombok</groupId>
            <artifactId>lombok</artifactId>
            <scope>provided</scope>
        </dependency>
    </dependencies>
</project>
`

// AdapterSchedulePOM is the adapter-schedule module POM template
const AdapterSchedulePOM = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <parent>
        <groupId>{{GROUP_ID}}</groupId>
        <artifactId>{{ARTIFACT_ID}}</artifactId>
        <version>{{VERSION}}</version>
        <relativePath>../../pom.xml</relativePath>
    </parent>

    <artifactId>adapter-schedule</artifactId>
    <packaging>jar</packaging>
    <name>adapter-schedule</name>
    <description>定时任务适配器</description>

    <dependencies>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>application-user</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter</artifactId>
        </dependency>
        <dependency>
            <groupId>org.projectlombok</groupId>
            <artifactId>lombok</artifactId>
            <scope>provided</scope>
        </dependency>
    </dependencies>
</project>
`

// ApplicationUserPOM is the application-user module POM template
const ApplicationUserPOM = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <parent>
        <groupId>{{GROUP_ID}}</groupId>
        <artifactId>{{ARTIFACT_ID}}</artifactId>
        <version>{{VERSION}}</version>
        <relativePath>../../pom.xml</relativePath>
    </parent>

    <artifactId>application-user</artifactId>
    <packaging>jar</packaging>
    <name>application-user</name>
    <description>用户业务应用层</description>

    <dependencies>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>domain</artifactId>
        </dependency>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>infrastructure</artifactId>
        </dependency>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>common</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter</artifactId>
        </dependency>
        <dependency>
            <groupId>org.mapstruct</groupId>
            <artifactId>mapstruct</artifactId>
        </dependency>
        <dependency>
            <groupId>org.projectlombok</groupId>
            <artifactId>lombok</artifactId>
            <scope>provided</scope>
        </dependency>
    </dependencies>
</project>
`

// StarterPOM is the starter module POM template
const StarterPOM = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <parent>
        <groupId>{{GROUP_ID}}</groupId>
        <artifactId>{{ARTIFACT_ID}}</artifactId>
        <version>{{VERSION}}</version>
    </parent>

    <artifactId>starter</artifactId>
    <packaging>jar</packaging>
    <name>starter</name>
    <description>启动模块</description>

    <dependencies>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>adapter-rest</artifactId>
        </dependency>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>adapter-schedule</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-actuator</artifactId>
        </dependency>
        <dependency>
            <groupId>io.micrometer</groupId>
            <artifactId>micrometer-registry-prometheus</artifactId>
        </dependency>
    </dependencies>

    <build>
        <plugins>
            <plugin>
                <groupId>org.springframework.boot</groupId>
                <artifactId>spring-boot-maven-plugin</artifactId>
                <executions>
                    <execution>
                        <goals>
                            <goal>repackage</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
        </plugins>
    </build>
</project>
`

// ApplicationModulePOM is the template for new application modules
const ApplicationModulePOM = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <parent>
        <groupId>{{GROUP_ID}}</groupId>
        <artifactId>{{ARTIFACT_ID}}</artifactId>
        <version>{{VERSION}}</version>
        <relativePath>../../pom.xml</relativePath>
    </parent>

    <artifactId>application-{{MODULE_NAME}}</artifactId>
    <packaging>jar</packaging>
    <name>application-{{MODULE_NAME}}</name>
    <description>{{MODULE_DESCRIPTION}}业务应用层</description>

    <dependencies>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>domain</artifactId>
        </dependency>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>infrastructure</artifactId>
        </dependency>
        <dependency>
            <groupId>{{GROUP_ID}}</groupId>
            <artifactId>common</artifactId>
        </dependency>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter</artifactId>
        </dependency>
        <dependency>
            <groupId>org.mapstruct</groupId>
            <artifactId>mapstruct</artifactId>
        </dependency>
        <dependency>
            <groupId>org.projectlombok</groupId>
            <artifactId>lombok</artifactId>
            <scope>provided</scope>
        </dependency>
    </dependencies>
</project>
`
