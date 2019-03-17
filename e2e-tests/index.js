var chakram = require('chakram'), expect = chakram.expect;
var compose = require('@terascope/docker-compose-js')("../docker-compose.yml");
var path = require('path');

const timeout = (ms) => new Promise(resolve => setTimeout(resolve, ms));

before(async () => {
    // bring up containers
    await compose.down({ "-v": null }); // also removes volumes
    await compose.up();
    await timeout(10 * 1000); // wait for db init
    console.log(await compose.ps());
})

after(async () => {
    // bring down containers
    await compose.down();
})

describe("/ping", () => {
    it("should respond with \"pong\"", async () => {
        const res = await chakram.get("http://localhost/ping");
        expect(res.error).to.be.null;
        expect(res).to.have.status(200);
        expect(res.body).to.equal("pong");
    });
});

describe("/register", () => {
    it("should register valid user", async () => {
        const res = await chakram.post("http://localhost/register", {
            email: "foo@bar.com",
            password: "password",
        })

        expect(res.error).to.be.null;
        expect(res).to.have.status(200);
    })

    it("should not register user with identical email address", async () => {
        const res = await chakram.post("http://localhost/register", {
            email: "foo@bar.com",
            password: "different_password",
        })

        expect(res.error).to.be.null;
        expect(res).to.have.status(409);
        expect(res).to.comprise.of.json({
            error: 'A user with email "foo@bar.com" already exists',
            code: 6,
            status: 409
        });
    })

    it("should not register user with invalid email address", async () => {
        const res = await chakram.post("http://localhost/register", {
            email: "invalid_email_address",
            password: "password",
        })

        expect(res.error).to.be.null;
        expect(res).to.have.status(400);
        expect(res).to.comprise.of.json({
            error: 'Email \"invalid_email_address\" is invalid',
            code: 3,
            status: 400
        });
    })

    it("should not register user with a password < 8 characters", async () => {
        const res = await chakram.post("http://localhost/register", {
            email: "bar@foo.com",
            password: "pass",
        })

        expect(res.error).to.be.null;
        expect(res).to.have.status(400);
        expect(res).to.comprise.of.json({
            error: 'Password is less than 8 characters',
            code: 3,
            status: 400
        });
    })
})
