########################################################################################
# Authentication endpoint and save as a session cookie in $webSession object

# Define the login URL and form data
$loginUrl = "http://127.0.0.1:5000/loginEndpoint"
$username = "j.doe@fonteyn-vakantieparken.nl"
$password = "Admin01!"

# Create the form data
$formData = @{
    username = $username
    password = $password
}

# Create a web session object to store cookies
$webSession = New-Object Microsoft.PowerShell.Commands.WebRequestSession

# Send the POST request to the login endpoint
$response = Invoke-WebRequest -Uri $loginUrl -Method Post -Body $formData -WebSession $webSession

########################################################################################
# Get all employees

# Check the response status code
if ($response.StatusCode -eq 200) {
    Write-Output "Successfully authenticated user"
    
    # The cookie is automatically stored in the $webSession object
    # You can use this $webSession object for subsequent requests
    Write-Output "Session cookies saved for future requests"
    

    $apiUrl = "http://127.0.0.1:5000/api/v1/employee"
    try {
        $response = Invoke-RestMethod -Uri $apiUrl -Method Get -WebSession $webSession
        Write-Output "API response ontvangen:"
        Write-Output $response
    } catch {
        Write-Output "Er is een fout opgetreden:"
        Write-Output $_.Exception.Message
    }
} else {
    Write-Output "Failed to authenticate user"
    Write-Output $response.Content
}

########################################################################################
