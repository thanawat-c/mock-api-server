// Complete ECMAScript Set methods polyfill for Node.js < 22
if (!Set.prototype.difference) {
  Set.prototype.difference = function (other) {
    const result = new Set();
    for (const elem of this) {
      if (!other.has(elem)) result.add(elem);
    }
    return result;
  };
}
if (!Set.prototype.intersection) {
  Set.prototype.intersection = function (other) {
    const result = new Set();
    for (const elem of this) {
      if (other.has(elem)) result.add(elem);
    }
    return result;
  };
}
if (!Set.prototype.union) {
  Set.prototype.union = function (other) {
    const result = new Set(this);
    for (const elem of other) result.add(elem);
    return result;
  };
}
if (!Set.prototype.isDisjointFrom) {
  Set.prototype.isDisjointFrom = function (other) {
    for (const elem of this) {
      if (other.has(elem)) return false;
    }
    return true;
  };
}
if (!Set.prototype.isSubsetOf) {
  Set.prototype.isSubsetOf = function (other) {
    for (const elem of this) {
      if (!other.has(elem)) return false;
    }
    return true;
  };
}
if (!Set.prototype.isSupersetOf) {
  Set.prototype.isSupersetOf = function (other) {
    for (const elem of other) {
      if (!this.has(elem)) return false;
    }
    return true;
  };
}
if (!Set.prototype.symmetricDifference) {
  Set.prototype.symmetricDifference = function (other) {
    const result = new Set(this);
    for (const elem of other) {
      if (this.has(elem)) result.delete(elem);
      else result.add(elem);
    }
    return result;
  };
}
