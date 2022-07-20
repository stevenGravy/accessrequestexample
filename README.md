# Auto Request Example
Provides an example using the Teleport API to submit and approve access requests.  
The program allows you to request multiple roles for a userid.  If the user has not logged in
via SSO or been created yet the program will wait by default until that userid exists.  

# Usage Example

```bash
$ accessrequestauto --proxy=enterprise.teleportdemo.com:443 --identity ./example-request.pem -user bob@example.dev -roles access                    

# 2022/07/19 21:06:27 Submitting access request for enterprise.teleportdemo.com:443 on bob@example.dev for roles [access]
# 2022/07/19 21:06:27 Connected to server: {tele1c 10.0.1 Kubernetes:true App:true DB:true OIDC:true SAML:true AccessControls:true AdvancedAccessWorkflows:true HSM:true Desktop:true ModeratedSessions:true MachineID:true ResourceAccessRequests:true  enterprise.teleportdemo.com:443 %!s(bool=false) {}  %!s(int32=0)}
# 2022/07/19 21:06:27 Access request state: APPROVED
```

The `bob@example.dev` user now has a 

# Prerequisites

- Existing Teleport Cluster
- Role assignable to a user that has a role request configuration

```yaml
kind: role
version: v5
metadata:
  name: devops
spec:
  allow:
    request:
      roles: ['access']
```
- Use the identity example below or a user with the same role rights as in `examplerole.yaml`

# Options

Run `accessrequestauto` without options to see all the parameters.  Using `--variable=true|false` for boolean parameters.

# Build the program

```bash
# go verion 1.18
go build accessrequestauto.go
```

# Configuring for usage

## Insert role and user

The example role and user has the right to create an access request
and approve any role.  Modify for further restriction.

### Insert role and user
```bash
tctl create -f examplerole.yaml
tctl create -f user.yaml
```

### Sign the user identity file

Note if you are using remotely the role will need the impersonate rights for your user.

```bash
tctl auth sign --user=example-request --out=example-request.pem --ttl=10000h
```



