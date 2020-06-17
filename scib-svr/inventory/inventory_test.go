package inventory

import (
	"net/http"
	"reflect"
	"scib-svr/datastore"
	"testing"
)

type InventoryServiceTest struct {
	service *Service
}

func NewInventoryServiceTest() *InventoryServiceTest {
	return &InventoryServiceTest{
		NewService(datastore.NewFirestoreDatastore()),
	}
}

const FirestoreEmulatorHost = "FIRESTORE_EMULATOR_HOST"

/*func TestMain(m *testing.M) {
	// start emulator
	cmd := exec.Command("gcloud", "beta", "emulators", "firestore", "start")
	// make it killable
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// we need to capture output to know when it's fully started up
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stderr.Close()
	// start
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// kill process in any case after finish tests
	var result int
	defer func() {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		os.Exit(result)
	}()
	// wait until it's running
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// read output
		buf := make([]byte, 256, 256)
		for {
			n, err := stderr.Read(buf[:])
			if err != nil {
				// until it ends
				if err == io.EOF {
					break
				}
				log.Fatalf("reading stderr %v", err)
			}
			if n > 0 {
				d := string(buf[:])
				log.Printf("%s", d)
				// has it started yet? Dad! Are we there ?
				if strings.Contains(d, "Dev App Server is now running") {
					wg.Done()
				}
				// capture the FIRESTORE_EMULATOR_HOST value to set
				pos := strings.Index(d, FirestoreEmulatorHost+"=")
				if pos > 0 {
					port := d[pos+len(FirestoreEmulatorHost)+5 : pos+len(FirestoreEmulatorHost)+9]
					err = os.Setenv(FirestoreEmulatorHost, fmt.Sprintf("localhost:%s", port))
					if err != nil {
						log.Fatalf("setting env failed %v:%v", err, port)
					}
				}
			}
		}
	}()
	wg.Wait()
	result = m.Run()
}*/

func TestCRUD(t *testing.T) {
	test := NewInventoryServiceTest()
	item := NewItem()
	item.Color = "navy blue"
	item.Size = "XL"
	item.Brand = "Swimma"
	item.Categories = []string{"Swimwear", "Sport", "Clothing"}
	item.Description = "Navy Blue Swimsuit for women"
	item.Name = "Swimma One WC"
	item.Sku = "swi-one-wc-nav-blue-xl"
	item.Cnt = 15
	item.Available = true
	item.Stockable = true
	status, id, err := test.service.save(item)
	if status != http.StatusCreated {
		t.Errorf("save failed because of wrong status: expected: %v got %v", http.StatusCreated, status)
	}
	if err != nil {
		t.Errorf("SAVE failed because of an error %v", err)
	}
	item.Id = id
	readItem, err := test.service.itemById(id)
	if err != nil {
		t.Errorf("READ failed because of an error %v", err)
	}
	if !reflect.DeepEqual(item, readItem) {
		t.Errorf("READ failed expected: %v got %v", item, readItem)
	}
	item.Name = "Swimma Two AC"
	status, id, err = test.service.save(item)
	if status != http.StatusNoContent {
		t.Errorf("save failed because of wrong status: expected: %v got %v", http.StatusNoContent, status)
	}
	if err != nil {
		t.Errorf("SAVE failed because of an error %v", err)
	}
	if reflect.DeepEqual(item, readItem) {
		t.Errorf("READ failed expected: %v got %v", item, readItem)
	}
	readItem, err = test.service.itemById(id)
	if err != nil {
		t.Errorf("READ failed because of an error %v", err)
	}
	if !reflect.DeepEqual(item, readItem) {
		t.Errorf("READ failed expected: %v got %v", item, readItem)
	}
	deleted, err := test.service.delete(id)
	if err != nil {
		t.Errorf("READ failed because of an error %v", err)
	}
	if !reflect.DeepEqual(deleted, readItem) {
		t.Errorf("READ failed expected: %v got %v", item, readItem)
	}
	readItem, err = test.service.itemById(id)
	if err == nil {
		t.Errorf("DELETE failed expected: %v got %v", nil, readItem)
	}
}


