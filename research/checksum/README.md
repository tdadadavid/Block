# Checksum

---

Checksum, as I understand it, is a way to verify the integrity of a file. I break down into two parts:
"Check" and "Sum." The sum is the value that is used to verify the integrity of the file, which is why it is called 'cheksum' because the 
sum is used to 'check' the file. There are many different types of checksums, but the most common are MD5, SHA1, and SHA256.
If you perform a checksum on a file using two different algorithms, the two checksums will be different because they all use different mathematical algorithms

### Common use cases
- File integrity
- Cryptographic verification
- Data transfer (i.e Downloading a file, uploading a file)

---

### Types of checksums
- XOR ✅
- CRC32 ✅
- SHA256 

_*NOTE*_: There are more, but I am stopping here.

### XOR
This takes two inputs and performs the bitwise XOR operation on them. 

```markdown
    a ^ a = 0 → a number XOR itself is always 0
    a ^ b = 1 → a number XOR a different is always 1
    b ^ 0 = 1 → a number XOR 0 is always 1
```

### CRC32 (Cyclic Redundancy Check)
This checksum is used to verify the integrity of a file. It is an error detection code. It is used where
data integrity is important but not cryptographic security is not required. It is used in many different
applications e.g, PNG, Ethernet, and ZIP files.
It uses a fixed polynomial (0xEDB88320) and a fixed initial value (0xFFFFFFFF).
- [Wikipedia](https://en.wikipedia.org/wiki/Cyclic_redundancy_check)

### References
- [Wikipedia](https://en.wikipedia.org/wiki/Checksum)
