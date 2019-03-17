var chakram = require('chakram'), expect = chakram.expect;
var compose = require('@terascope/docker-compose-js')("../docker-compose.yml");
var path = require('path');

before(async () => {
    // bring up containers
    await compose.down();
    await compose.rm(null, { "-v": null }); // remove volumes
    await compose.up();
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