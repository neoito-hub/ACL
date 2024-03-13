## ACL Service

### Overview

Welcome to the ACL Service for AppBlocks, a comprehensive solution for managing access control within your application environment. This service comprises several essential components designed to ensure robust authentication and authorization mechanisms. Let's delve into the key components and steps for seamless integration.

### Components

#### 1. Shield

Shield plays a pivotal role in authenticating applications within your ecosystem. It serves as the gateway for validating access requests, ensuring only authorized entities interact with your system.

#### 2. Captain

Captain orchestrates the management of APIs, acting as a central hub for routing and controlling access to various endpoints. It works hand in hand with Spaces to enforce access policies effectively.

#### 3. Spaces

Spaces complements Shield by providing a framework for managing access policies. It allows granular control over who can access specific resources and endpoints mapped through Captain.

#### 4. Data Models

In this section, we outline our data models and migrations.

#### 5. Frontend Application

The frontend application provides users with a user-friendly interface to interact with your system. It communicates with the backend components, leveraging the authentication and authorization mechanisms provided by Shield and Captain.

# Integration Guide for ACL Block

This guide outlines the steps necessary to integrate the ACL (Access Control List) block, a repository that includes Docker Compose configurations and code necessary for implementing access control mechanisms, into your application.

## Prerequisites

- Docker and Docker Compose installed on your machine.
- Basic understanding of Docker, environment variables, and access control mechanisms.

## Step 1: Clone the ACL Block Repository

Clone the ACL block repository to your local machine. This repository, referred to as a "block," contains all necessary Docker Compose configurations and access control code.

```bash
git clone <repository-url>
```

_Replace `<repository-url>` with the actual URL of the ACL block repository._

## Step 2: Configure Environment Variables

1. Navigate to the repository directory after cloning.
2. Locate the `.env.sample` files.
3. Copy each `.env.sample` to a new `.env` file.
4. Edit the `.env` files, replacing placeholder values with your specific configurations.

## Step 3: Build and Run the Application with Docker Compose

With the environment variables set, build and start the application services using Docker Compose:

```bash
docker compose build && docker compose up
```

## Step 4: Seed Data for ACL

Seed the application with necessary resources, actions, policies, and permissions. Refer to `DataModels/migrator/main.go` for an example on how to seed data.

## Step 5: Integrate ACL in the Browser Application

To integrate ACL in a browser application:

1. Integrate `@ppblocks/js-sdk`.
2. Set the environment variable for the Shield authentication URL:

```bash
BB_SHIELD_AUTH_URL=http://localhost:8011/
```

## Step 6: Add ACL Middleware to Server Application

Add middleware to your server application that checks user access against the ACL endpoint and responds accordingly. Use the following Axios request as a template:

```javascript
await axios({
  method: "get",
  url: `http://localhost:6001/api/auth/getUser?action-name=<action>&space-id=<your space id>&resource-name=<resource name>`,
  headers: {
    accept: "application/json",
    "content-type": "application/json",
    authorization: `<auth token from shield>`,
  },
});
```

\*Replace <action>, <your space id>, <resource name>, and <auth token from shield> with the appropriate values for your application.

### Implementation Details

#### 1. App Seeding

To utilize the Shield service, we must include our application details in the `shield_apps table`, which is linked with the bridge table named `shield_app_domain_mappings`. Additionally, the application's required permissions are stored in the `permissions` table. Furthermore, there exists a bridge table named `app_permissions` that connects the applications with their corresponding permissions.

#### 2. Entity Configuration

We are seeding entity types and space access entities . An entity refers to the attachment of users to an entity, each consisting of specific policies. These entities are categorized under entity types. We maintain a table (`entity_type_definitions`) to store definitions of entity types. Additionally, we seed space access entities into the (`entities`) table.

#### 3. Resource Initialization

To begin, we need to populate the necessary tables with resource data. This includes adding API endpoints to the `ac_resource` table. Additionally, we utilize secondary tables such as `ac_res_grps` and a bridge table named `ac_res_gp_res` for organizing and linking resources efficiently.

#### 4. Policy Definition

Next, we define access policies by adding entries to the `ac_policies` table. Similar to resources, we utilize secondary tables (`ac_pol_grps`) and bridge tables (`pol_gp_policies`) to structure and manage policies effectively.

#### 5. Action Management

We maintain an `ac_actions` table to store API action types and utilize bridge tables (`ac_act_grps`, `act_gp_actions`) for associating actions with groups and resources, facilitating fine-grained access control.

#### 6. Permission Configuration

Permissions are crucial for defining the level of access granted to users or roles. We utilize the `ac_permissions` table to store permission details, with bridge tables (`per_pol_grps`) facilitating the association between permissions and policies.

### Integration Guidelines

When integrating with the ACL service, ensure to follow these guidelines:

- Utilize migration scripts provided within the Data Models for seeding initial data.
- When registering a new application, use the migration query while replacing the ID and secret parameters as required.
- Maintain consistency in managing resources, policies, actions, and permissions to ensure a robust access control setup.

By adhering to these guidelines and leveraging the capabilities of the ACL service, you can enhance the security and manageability of your application ecosystem effectively.

### Getting Started

To get started with the ACL service:

1. Install Docker on your machine.
2. Navigate to the root directory of the project.
3. Run `docker-compose up` to deploy the necessary containers and start the service.

With the ACL service up and running, you can begin configuring access control policies and managing resources to suit your application's requirements.

### Conclusion

The ACL service for AppBlocks provides a comprehensive solution for managing access control within your application environment. By understanding its components and following the implementation steps and guidelines provided, you can ensure secure and efficient access management for your users and applications.
