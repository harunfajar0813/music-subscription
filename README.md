# Music Subscription Submission

![meme-programmer](https://cupheadmemes.com/wp-content/uploads/2018/08/Best-Programming-Memes-80.jpg) 

### DON'T FORGET TO FOLLOW 7 DEADLY CONVENTION, MENTIONED IN THIS REPOSITORIES ON [CONVENTION](https://github.com/Aldiwildan77/music-subscription/blob/master/CONVENTION.md)
### Any code that didn't follow the convention will automatically rejected

## <b>References</b> 
1. [Example of Project](https://github.com/meong1234/fintech)
2. [Git](https://try.github.io/)
3. [CheatSheets](https://devhints.io/)
4. [REST API](https://restfulapi.net/)
5. [Insomnia REST Client](https://insomnia.rest/)
6. [Test Driven Development](https://www.freecodecamp.org/news/test-driven-development-what-it-is-and-what-it-is-not-41fa6bca02a2/)

## <b>Accepted Weapon</b>
1. NodeJS (Javascript)
2. Golang
3. Java

## <b>Problem Statement</b>
One day you are thinking of building a music subscription apps that can do non-cash transactions, as the `Minimum Viable Product (MVP)` you want this application to have the following these features:

* A new customer can register to the system
* Customer can do top-up balance
* Customer can buy a subscription
* Customer can renew a subscription

## <b>Meet the actor</b>
```
1. Customer
2. Subscription 
3. Transaction
```

## <b>Some Spec</b>
```
GIVEN I am unregistered person
WHEN I register as customer with (name, email, phone)
THEN Customer should be record as new customer and return Customer ID

GIVEN I am Customer
WHEN I top-up some amount balance
THEN Customer balance should be increased

GIVEN I am Customer 
WHEN I buy a subscription
THEN Customer balance should be decreased, Transaction recorded as receipt, and return Transaction ID

GIVEN I am Customer
WHEN I need to renew a subscription
THEN Customer balance should be decreased, Transaction recorded as receipt, and return Transaction ID
```

## <b>Entities</b>
* Customer: 
  * id
  * name
  * email
  * phone
  * balance

* Subscription 
  * id
  * name
  * price
  * duration

* Transaction
  * id
  * customer_id
  * subscription_id
  * total

## <b>API</b>
```
-> /customer/register
  {
    "name": ""
    "email": ""
    "phone": ""
  }
-> /customer/topup
  {
    "customer_id": ""
    "amount": ""
  }
-> /transaction/payment
  {
    "customer_id": ""
    "subscription_id": ""
    "total": ""
  }
```
> ### You can add a new API for exploration but you `must` implement the API above

## <b>TODO</b>
* [ ] Setup your Environment of this project by your chosen language
* [ ] Prepare the actors services (Controller, Model, Route)
  * [ ] Customer
    * [ ] Register
    * [ ] Topup
    * [ ] Debit
  * [ ] Subscription
    * [ ] Create Subscription
    * [ ] Read Subscription
    * [ ] Read Subscription By Id
  * [ ] Transaction
    * [ ] Create Transaction
    * [ ] Read Transaction 
    * [ ] Read Transaction By ID
* [ ] Do Testing

## <b>TASK</b>
```
1. Fork This Repository
2. Follow the project convention
3. Finish all TODO
4. Implement Test Driven Development
```

## <b>NOTES</b>
``` 
1. `Folder structure` can be modified by your own  
2. Implement `REST API` will be a score plus
3. ONLY BUSINESS LOGIC (BACKEND) 
4. Don't forget to update your .gitignore file
5. Code with your soulmate (Pairing) or do it alone
```

## <b>HOW TO SUBMIT</b>
Please follow the [CONTRIBUTING](https://github.com/Aldiwildan77/music-subscription/blob/master/CONTRIBUTING.md)

## This is not the only way to join with us
## But this is the one and only way to instantly pass
