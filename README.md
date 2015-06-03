# sphere-thing-deleter

Command line tool for deleting Ninja Sphere "things"

**usage:** sphere-thing-deleter [method] [value]


Supported methods: `type`, `name`, `promoted`

### Examples:

 - To delete all non-promoted things, use: `                ... promoted false`
 - To delete all things with type 'light', use: `           ... type light`
 - To delete all things with names containing 'jim', use: ` ... name jim`