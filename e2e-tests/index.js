var chakram = require('chakram');

describe("Chakram", function() {
    it("should offer simple HTTP request capabilities", function () {
        return chakram.get("http://httpbin.org/get");
    });
});