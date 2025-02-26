# Credentials Folder

## The purpose of this folder is to store all credentials needed to log into your server and databases. This is important for many reasons. But the two most important reasons is
    1. Grading , servers and databases will be logged into to check code and functionality of application. Not changes will be unless directed and coordinated with the team.
    2. Help. If a class TA or class CTO needs to help a team with an issue, this folder will help facilitate this giving the TA or CTO all needed info AND instructions for logging into your team's server. 


# Below is a list of items required. Missing items will causes points to be deducted from multiple milestone submissions.

1. 54.151.68.50:8081 (This is also the link for website)
2. unbuntu
3. Using key
4. Database URL or IP and port used.
5. Database username
6. Database password
7. Database name (basically the name that contains all your tables)
8. Instructions on how to use the above information.

## How to Access and Manage the System

> **1. Download the Key**
> - Obtain the SSH key file (`CSC848.pem`) from a secure source.

> **2. Navigate to Your Downloads Folder**
> - Open a terminal and run:
>   ```bash
>   cd "Your download path here"
>   ```
> - Replace `"Your download path here"` with the actual path where your SSH key is located.

> **3. Connect to the Server**
> - Use SSH to log into the AWS EC2 instance:
>   ```bash
>   ssh -i "CSC848.pem" ubuntu@ec2-54-151-68-50.us-west-1.compute.amazonaws.com
>   ```
> - Ensure that the `.pem` file has the correct permissions:
>   ```bash
>   chmod 400 CSC848.pem
>   ```

> **4. Navigate to the Application Directory**
> - Once logged into the EC2 instance, run:
>   ```bash
>   cd csc<TAB>
>   ```
> - Press **Tab** to autocomplete the folder name if needed.

> **5. Start and Stop the Server**
> - **To stop the server:**
>   ```bash
>   systemctl --user stop GoApp.service
>   ```
> - **To start the server:**
>   ```bash
>   systemctl --user start GoApp.service
>   

# Most important things to Remember
## These values need to kept update to date throughout the semester. <br>
## <strong>Failure to do so will result it points be deducted from milestone submissions.</strong><br>
## You may store the most of the above in this README.md file. DO NOT Store the SSH key or any keys in this README.md file.
