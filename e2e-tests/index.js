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

describe("Chakram", function() {
    it("should offer simple HTTP request capabilities", function () {
        return chakram.get("http://httpbin.org/get");
    });
});