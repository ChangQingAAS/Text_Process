## golang os.OpenFile几种常用模式
- os.O_WRONLY | os.O_CREATE | O_EXCL           【如果已经存在，则失败】
- os.O_WRONLY | os.O_CREATE                         【如果已经存在，会覆盖写，不会清空原来的文件，而是从头直接覆盖写】
- os.O_WRONLY | os.O_CREATE | os.O_APPEND  【如果已经存在，则在尾部添加写】

一般都文件属性标识如下：

-rwxrwxrwx

第1位：文件属性，一般常用的是"-"，表示是普通文件；"d"表示是一个目录。

第2～4位：文件所有者的权限rwx (可读/可写/可执行)。

第5～7位：文件所属用户组的权限rwx (可读/可写/可执行)。

第8～10位：其他人的权限rwx (可读/可写/可执行)。



在golang中，可以使用os.FileMode(perm).String()来查看权限标识：

os.FileMode(0777).String()    //返回 -rwxrwxrwx

os.FileMode(0666).String()   //返回 -rw-rw-rw-

os.FileMode(0644).String()   //返回 -rw-r--r--



0777表示：创建了一个普通文件，所有人拥有所有的读、写、执行权限

0666表示：创建了一个普通文件，所有人拥有对该文件的读、写权限，但是都不可执行

0644表示：创建了一个普通文件，文件所有者对该文件有读写权限，用户组和其他人只有读权限，都没有执行权限

注意，golang中创建文件指定权限时，只能以"0XXX"的形式，不能省掉前面的"0"，否则指定的权限不是预期的。如：

os.FileMode(777).String()   //返回 -r----x--x

os.FileMode(666).String()   //返回 --w--wx-w-

os.FileMode(644).String()   //返回 --w----r--