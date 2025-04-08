package storage

import "strings"

// HashMap方法
func NewHashMap(capacity int, hashFunc HashFunc, compareFunc CompareFunc) *HashMap {
    buckets := make([]Bucket, capacity)
    for i := range buckets {
        buckets[i] = NewLinkedList(compareFunc)
    }
    return &HashMap{
        hashFunc:    hashFunc,
        compareFunc: compareFunc,
        buckets:     buckets,
        capacity:    capacity,
    }
}

func (hm *HashMap) Put(key, value string) {
    idx := hm.hashFunc(key) % hm.capacity
    bucket := hm.buckets[idx]

    inserted := bucket.Put(key, value)
    if inserted {
        hm.size++
    }

    if ll, ok := bucket.(*LinkedList); ok && ll.Len() >= 8 {
        rb := convertLinkedListToRBTree(ll, hm.compareFunc)
        hm.buckets[idx] = rb
    }
}

func (hm *HashMap) Get(key string) (string, bool) {
    idx := hm.hashFunc(key) % hm.capacity
    value, b := hm.buckets[idx].Get(key)
	return value.(string), b
}

func (hm *HashMap) Del(key string) {
    idx := hm.hashFunc(key) % hm.capacity
    bucket := hm.buckets[idx]

    if bucket.Delete(key) {
        hm.size--
        if rb, ok := bucket.(*RBTree); ok && rb.Len() <= 6 {
            ll := convertRBTreeToLinkedList(rb, hm.compareFunc)
            hm.buckets[idx] = ll
        }
    }
}

func (hm *HashMap) Prefix(prefix string) []string {
    result := make([]string, 0)
    entries := hm.Entries()
    for _, entry := range entries {
        key := entry.Key.(string)
        if key == "" {
			continue
		}
        if strings.HasPrefix(key, prefix) {
            result = append(result, key)
        }
    }
    return result
}

func (hm *HashMap) Suffix(suffix string) []string {
    result := make([]string, 0)
    entries := hm.Entries()
    for _, entry := range entries {
        key := entry.Key.(string)
        if key == "" {
			continue
		}
        if strings.HasSuffix(key, suffix) {
            result = append(result, key)
        }
    }
    return result
}

func (hm *HashMap) Contains(sub string) []string {
    result := make([]string, 0)
    entries := hm.Entries()
    for _, entry := range entries {
        key := entry.Key.(string)
        if key == "" {
			continue
		}
        if strings.HasSuffix(key, sub) {
            result = append(result, key)
        }
    }
    return result
}

// 新增Entry结构体，用于存储键值对
type Entry struct {
    Key   interface{}
    Value interface{}
}

// HashMap新增Entries方法
func (hm *HashMap) Entries() []Entry {
    entries := make([]Entry, 0, hm.size)
    for _, bucket := range hm.buckets {
        switch b := bucket.(type) {
        case *LinkedList:
            curr := b.head
            for curr != nil {
                entries = append(entries, Entry{Key: curr.key, Value: curr.value})
                curr = curr.next
            }
        case *RBTree:
            var traverse func(*RBTreeNode)
            traverse = func(node *RBTreeNode) {
                if node == nil {
                    return
                }
                traverse(node.left)
                entries = append(entries, Entry{Key: node.key, Value: node.value})
                traverse(node.right)
            }
            traverse(b.root)
        }
    }
    return entries
}

type Color bool

const (
    Red   Color = true
    Black Color = false
)

type CompareFunc func(a, b interface{}) int
type HashFunc func(key interface{}) int

type HashMap struct {
    hashFunc    HashFunc
    compareFunc CompareFunc
    buckets     []Bucket
    size        int
    capacity    int
}

type Bucket interface {
    Put(key, value interface{}) bool
    Get(key interface{}) (interface{}, bool)
    Delete(key interface{}) bool
    Len() int
}

// 链表实现
type LinkedList struct {
    head        *ListNode
    length      int
    compareFunc CompareFunc
}

type ListNode struct {
    key   interface{}
    value interface{}
    next  *ListNode
}

func NewLinkedList(compareFunc CompareFunc) *LinkedList {
    return &LinkedList{compareFunc: compareFunc}
}

func (ll *LinkedList) Put(key, value interface{}) bool {
    curr := ll.head
    for curr != nil {
        if ll.compareFunc(curr.key, key) == 0 {
            curr.value = value
            return false
        }
        curr = curr.next
    }
    ll.head = &ListNode{key, value, ll.head}
    ll.length++
    return true
}

func (ll *LinkedList) Get(key interface{}) (interface{}, bool) {
    curr := ll.head
    for curr != nil {
        if ll.compareFunc(curr.key, key) == 0 {
            return curr.value, true
        }
        curr = curr.next
    }
    return nil, false
}

func (ll *LinkedList) Delete(key interface{}) bool {
    var prev *ListNode
    curr := ll.head
    for curr != nil {
        if ll.compareFunc(curr.key, key) == 0 {
            if prev == nil {
                ll.head = curr.next
            } else {
                prev.next = curr.next
            }
            ll.length--
            return true
        }
        prev = curr
        curr = curr.next
    }
    return false
}

func (ll *LinkedList) Len() int { return ll.length }

// 红黑树实现
type RBTree struct {
    root        *RBTreeNode
    compareFunc CompareFunc
    length      int
}

type RBTreeNode struct {
    key    interface{}
    value  interface{}
    color  Color
    left   *RBTreeNode
    right  *RBTreeNode
    parent *RBTreeNode
}

func NewRBTree(compareFunc CompareFunc) *RBTree {
    return &RBTree{compareFunc: compareFunc}
}

func (rb *RBTree) Put(key, value interface{}) bool {
    newNode := &RBTreeNode{key: key, value: value, color: Red}
    if rb.root == nil {
        newNode.color = Black
        rb.root = newNode
        rb.length = 1
        return true
    }

    current := rb.root
    var parent *RBTreeNode
    for current != nil {
        parent = current
        cmp := rb.compareFunc(key, current.key)
        switch {
        case cmp < 0:
            current = current.left
        case cmp > 0:
            current = current.right
        default:
            current.value = value
            return false
        }
    }

    newNode.parent = parent
    cmp := rb.compareFunc(key, parent.key)
    if cmp < 0 {
        parent.left = newNode
    } else {
        parent.right = newNode
    }

    rb.insertFixup(newNode)
    rb.length++
    return true
}

func (rb *RBTree) insertFixup(x *RBTreeNode) {
    for x.parent != nil && x.parent.color == Red {
        if x.parent == x.parent.parent.left {
            y := x.parent.parent.right
            if y != nil && y.color == Red {
                x.parent.color = Black
                y.color = Black
                x.parent.parent.color = Red
                x = x.parent.parent
            } else {
                if x == x.parent.right {
                    x = x.parent
                    rb.rotateLeft(x)
                }
                x.parent.color = Black
                x.parent.parent.color = Red
                rb.rotateRight(x.parent.parent)
            }
        } else {
            y := x.parent.parent.left
            if y != nil && y.color == Red {
                x.parent.color = Black
                y.color = Black
                x.parent.parent.color = Red
                x = x.parent.parent
            } else {
                if x == x.parent.left {
                    x = x.parent
                    rb.rotateRight(x)
                }
                x.parent.color = Black
                x.parent.parent.color = Red
                rb.rotateLeft(x.parent.parent)
            }
        }
    }
    rb.root.color = Black
}

func (rb *RBTree) rotateLeft(x *RBTreeNode) {
    y := x.right
    x.right = y.left
    if y.left != nil {
        y.left.parent = x
    }
    y.parent = x.parent
    if x.parent == nil {
        rb.root = y
    } else if x == x.parent.left {
        x.parent.left = y
    } else {
        x.parent.right = y
    }
    y.left = x
    x.parent = y
}

func (rb *RBTree) rotateRight(y *RBTreeNode) {
    x := y.left
    y.left = x.right
    if x.right != nil {
        x.right.parent = y
    }
    x.parent = y.parent
    if y.parent == nil {
        rb.root = x
    } else if y == y.parent.right {
        y.parent.right = x
    } else {
        y.parent.left = x
    }
    x.right = y
    y.parent = x
}

func (rb *RBTree) Get(key interface{}) (interface{}, bool) {
    node := rb.findNode(key)
    if node != nil {
        return node.value, true
    }
    return nil, false
}

func (rb *RBTree) findNode(key interface{}) *RBTreeNode {
    current := rb.root
    for current != nil {
        cmp := rb.compareFunc(key, current.key)
        switch {
        case cmp < 0:
            current = current.left
        case cmp > 0:
            current = current.right
        default:
            return current
        }
    }
    return nil
}

func (rb *RBTree) Delete(key interface{}) bool {
    node := rb.findNode(key)
    if node == nil {
        return false
    }

    var child, parent *RBTreeNode
    color := node.color

    if node.left == nil {
        child = node.right
        rb.transplant(node, node.right)
        parent = node.parent
    } else if node.right == nil {
        child = node.left
        rb.transplant(node, node.left)
        parent = node.parent
    } else {
        successor := rb.minimum(node.right)
        color = successor.color
        child = successor.right
        parent = successor.parent

        if successor.parent != node {
            rb.transplant(successor, successor.right)
            successor.right = node.right
            successor.right.parent = successor
        }

        rb.transplant(node, successor)
        successor.left = node.left
        successor.left.parent = successor
        successor.color = node.color
    }

    if color == Black {
        rb.deleteFixup(child, parent)
    }

    rb.length--
    return true
}

func (rb *RBTree) transplant(u, v *RBTreeNode) {
    if u.parent == nil {
        rb.root = v
    } else if u == u.parent.left {
        u.parent.left = v
    } else {
        u.parent.right = v
    }
    if v != nil {
        v.parent = u.parent
    }
}

func (rb *RBTree) minimum(node *RBTreeNode) *RBTreeNode {
    for node.left != nil {
        node = node.left
    }
    return node
}

func (rb *RBTree) deleteFixup(x *RBTreeNode, parent *RBTreeNode) {
    for x != rb.root && (x == nil || x.color == Black) {
        if x == parent.left {
            sibling := parent.right
            if sibling.color == Red {
                sibling.color = Black
                parent.color = Red
                rb.rotateLeft(parent)
                sibling = parent.right
            }
            if (sibling.left == nil || sibling.left.color == Black) &&
               (sibling.right == nil || sibling.right.color == Black) {
                sibling.color = Red
                x = parent
                parent = x.parent
            } else {
                if sibling.right == nil || sibling.right.color == Black {
                    if sibling.left != nil {
                        sibling.left.color = Black
                    }
                    sibling.color = Red
                    rb.rotateRight(sibling)
                    sibling = parent.right
                }
                sibling.color = parent.color
                parent.color = Black
                if sibling.right != nil {
                    sibling.right.color = Black
                }
                rb.rotateLeft(parent)
                x = rb.root
            }
        } else {
            sibling := parent.left
            if sibling.color == Red {
                sibling.color = Black
                parent.color = Red
                rb.rotateRight(parent)
                sibling = parent.left
            }
            if (sibling.left == nil || sibling.left.color == Black) &&
               (sibling.right == nil || sibling.right.color == Black) {
                sibling.color = Red
                x = parent
                parent = x.parent
            } else {
                if sibling.left == nil || sibling.left.color == Black {
                    if sibling.right != nil {
                        sibling.right.color = Black
                    }
                    sibling.color = Red
                    rb.rotateLeft(sibling)
                    sibling = parent.left
                }
                sibling.color = parent.color
                parent.color = Black
                if sibling.left != nil {
                    sibling.left.color = Black
                }
                rb.rotateRight(parent)
                x = rb.root
            }
        }
    }
    if x != nil {
        x.color = Black
    }
}

func (rb *RBTree) Len() int { return rb.length }

func convertLinkedListToRBTree(ll *LinkedList, cf CompareFunc) *RBTree {
    rb := NewRBTree(cf)
    curr := ll.head
    for curr != nil {
        rb.Put(curr.key, curr.value)
        curr = curr.next
    }
    return rb
}

func convertRBTreeToLinkedList(rb *RBTree, cf CompareFunc) *LinkedList {
    ll := NewLinkedList(cf)
    var traverse func(*RBTreeNode)
    traverse = func(node *RBTreeNode) {
        if node == nil {
            return
        }
        traverse(node.left)
        ll.Put(node.key, node.value)
        traverse(node.right)
    }
    traverse(rb.root)
    return ll
}