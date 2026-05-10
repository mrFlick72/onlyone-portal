Using the Playwrite MCP I want to test the following scenario, ONLY TEST NO FIX HAS TO BE DONE it is a end o end ai agentic test.

- open the browser in incognito the infrastructure support http and if you open the regular browser may be that teh user will silently redirect form http to https so the login will fails.
    - use the command `google-chrome --incognito --remote-debugging-port=9222 -unsafely-treat-insecure-origin-as-secure=http://local.api.vauthenticator.com:9090,http://local.management.vauthenticator.com:8085,http://local.onlyone-portal.com:8070 --disable-features=HttpsUpgrades,HttpsFirstBalancedMode,HttpsFirstModeIncognito  --user-data-dir=/tmp/chrome-e2e`
- try to open the utl http://local.onlyone-portal.com:8070/index
- maybe you will be redirected to the login page of VAuthenticator it depends to already existing sessions
- let's use the credentials loaded in the .env file USER_NAME key should be for the username field while PWD for the password field