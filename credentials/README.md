# Credentials Folder

This folder contains all credentials and access instructions required to log into your team’s **server** and **database**. Keeping this information accurate and up to date is essential for both **grading** and **technical support**.

---

## Purpose

1. **Grading** – Instructors, TAs, and the CTO will use these credentials to verify your application’s deployment and functionality.  
2. **Support** – If assistance is needed, this folder allows the TA or CTO to access your system using clear instructions.

---

## Required Information

| Item                           | Description                                                    |
|--------------------------------|----------------------------------------------------------------|
| **Website / Server URL**       | `http://204.236.166.51:9081`                                   |
| **SSH Username**               | `ubuntu`                                                       |
| **SSH Authentication**         | Password-based login or pem                                    |
| **Database IP / URL**          | `<enter database IP or hostname>`                              |
| **Database Port**              | `<8081>`                                                       |
| **Database Username**          | `<enter database username>`                                    |
| **Database Password**          | `HorseMomDadHouseThing1!`                                      |
| **Database Name**              | `N\A`                                                          |

> 🔁 **Note:** These values must be kept up to date throughout the semester. Missing or incorrect information will result in point deductions during milestone evaluations.

---

## 🛠️ How to Access and Manage the System

### 1. **Log In to the Server**

> This project uses **password authentication** — or the `Csc648.pem` key.

```bash
ssh ubuntu@204.236.166.51
```

When prompted, enter the team-provided password.

---

### 2. **Navigate to the Application Directory**

Once logged in:

```bash
cd ~/csc648-fa25-0104-team03/application/Backend
```

---

### 3. **Start or Stop the Go Server**

To **start** the Go backend service:

```bash
sudo systemctl start GoApp.service
```

To **stop** the Go backend service:

```bash
sudo systemctl stop GoApp.service
```

To check its status:

```bash
sudo systemctl status GoApp.service
```

---

## Important Reminders

- Make sure this file always reflects the current setup.
- Keep this folder organized — it will be checked during grading.
