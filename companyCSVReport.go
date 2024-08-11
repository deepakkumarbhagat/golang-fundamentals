
// CompanyCsvReportHeader is the header fields of the company csv report.
const CompanyCsvReportHeader string = "id,name,department,manager"

const (
	Corporate string = "Corporate" // Corporate represents Corporate department
	Delimiter string = ","
)

// Employee is a representation of the employee details.
type Employee struct {
	ID        uint64  `json:"id"`
	Name      string  `json:"name"`
	ManagerID *uint64 `json:"manager_id,omitempty"`
}

// Department is a representation of the department details.
type Department struct {
	ID               uint64 `json:"id"`
	Name             string `json:"name"`
	DepartmentHeadID uint64 `json:"department_head_id"`
}

// Company is a representation of the company details.
// It compromises of list of employees and departments.
// The struct is expected to be unmarshalled from the json input.
type Company struct {
	Employees   []Employee   `json:"employees"`
	Departments []Department `json:"departments"`
}

// CompanyCsvReportItem represents an item in the company csv report.
type CompanyCsvReportItem struct {
	EmployeeID uint64
	Name       string
	Department string
	Manager    string
}

// SortReport sorts report which is a slice of CompanyCsvReportItem.
// The report is sorted in ascending order of department, ceo/boss, head of department, name.
func sortCompanyCsvReportItems(report []CompanyCsvReportItem, hodMap map[uint64]string) {
	sort.Slice(report, func(i, j int) bool {

		// Sort by department
		if report[i].Department < report[j].Department {
			return true
		}

		if report[i].Department == report[j].Department {
			// Sort by CEO/BOSS.
			// Empty Manager implies either CEO/BOSS
			if report[i].Manager == "" {
				return true
			}
			// Sort by department head
			// Check department map to find if EmployeeID is head of department.
			if _, ok := hodMap[report[i].EmployeeID]; ok {
				return true
			}
			// If report[j].Manager is CEO/BOSS, return false
			if report[j].Manager == "" {
				return false
			}
			// If report[j].EmployeeID is head of department, return false
			if _, ok := hodMap[report[j].EmployeeID]; ok {
				return false
			}
			// Sort by Name
			if report[i].Name < report[j].Name {
				return true
			}
		}

		return false
	})
}

func solution() []string {
	/*jsonDump := `{
		"employees":[
		   {
			  "id":1,
			  "name":"Mary",
			  "manager_id":null
		   },
		   {
			  "id":2,
			  "name":"Steve",
			  "manager_id":1
		   },
		   {
			  "id":3,
			  "name":"Lee",
			  "manager_id":2
		   },
		   {
			  "id":4,
			  "name":"Cindy",
			  "manager_id":3
		   }
		],
		"departments":[
		   {
			  "id":1,
			  "name":"HR",
			  "department_head_id":2
		   },
		   {
			  "id":2,
			  "name":"Recruiting",
			  "department_head_id":3
		   }
		]
	 }`


	*/
	jsonDump := `{
				"employees":[
				   {
					  "id":1,
					  "name":"Tom",
					  "manager_id":null
				   },
				   {
					  "id":2,
					  "name":"Jennifer",
					  "manager_id":1
				   },
				   {
					  "id":6,
					  "name":"Mary",
					  "manager_id":1
				   },
				   {
					  "id":3,
					  "name":"Jack",
					  "manager_id":6
				   },
				   {
					"id":4,
					"name":"Tiler",
					"manager_id":3
				 },
				 {
					"id":5,
					"name":"Elizabeth",
					"manager_id":3
				 },
				 {
					"id":7,
					"name":"Sophia",
					"manager_id":2
				 }
				],
				"departments":[
				   {
					  "id":1,
					  "name":"Sales",
					  "department_head_id":3
				   },
				   {
					  "id":2,
					  "name":"Development",
					  "department_head_id":2
				   }
				]
			 }`

	// Unmarshal jsonDump into companyData
	var companyData Company
	err := json.Unmarshal([]byte(jsonDump), &companyData)
	if err != nil {
		fmt.Printf("could not unmarshal json: %v\n", err)
		return []string{}
	}

	// Create head of department (hod) map
	hodMap := make(map[uint64]string)
	for _, department := range companyData.Departments {
		hodMap[department.DepartmentHeadID] = department.Name
	}

	employeeMap := make(map[uint64]string)
	//this may cause extra space usage, as this is used only for fetch the
	// manager name from manager id, may be we can use reports map only for this.
	for _, emp := range companyData.Employees {
		employeeMap[emp.ID] = emp.Name
	}

	companyCsvReportMap := make(map[uint64]CompanyCsvReportItem)

	// Create report item for each employee
	// Each report item consists of {id, name, department, manager}
	for _, employee := range companyData.Employees {
		var reportItem CompanyCsvReportItem

		// Copy id, name for the employee
		reportItem.EmployeeID = employee.ID
		reportItem.Name = employee.Name

		// Get manager name of the employee only if
		// employee manager is not nil
		if employee.ManagerID != nil {
			if manager, ok := employeeMap[*employee.ManagerID]; ok {
				reportItem.Manager = manager
			}
		}

		// Get department of the employee
		// If the employee is department head, get the department from the hodMap
		if department, ok := hodMap[employee.ID]; ok {
			reportItem.Department = department
		} else {
			// If the employee has not manager, assign the department as Corporate
			if employee.ManagerID == nil {
				reportItem.Department = Corporate
			} else {
				// Otherwise, the department is same as employee's manager department
				if repItem, ok := companyCsvReportMap[*employee.ManagerID]; ok {
					reportItem.Department = repItem.Department
				}
			}
		}
		companyCsvReportMap[reportItem.EmployeeID] = reportItem
	}

	var csvReportItems []CompanyCsvReportItem
	for _, csvReportItem := range companyCsvReportMap {
		csvReportItems = append(csvReportItems, csvReportItem)
	}

	sortCompanyCsvReportItems(csvReportItems, hodMap)

	companyCsvReport := generateCompanyCsvReport(csvReportItems)

	fmt.Println(companyCsvReport)

	return companyCsvReport
}

// generateCompanyCsvReport generates company's csv report in required format.
// It returns slice of string companyCsvReport from the given csvReportItems
// which is a slice of CompanyCsvReportItem
// Sample output below.
// [id,name,department,manager 1,Mary,Corporate, 2,Steve,HR,Mary 3,Lee,Recruiting,Steve 4,Cindy,Recruiting,Lee]
func generateCompanyCsvReport(csvReportItems []CompanyCsvReportItem) []string {
	companyCsvReport := []string{CompanyCsvReportHeader}

	for _, csvReportItem := range csvReportItems {

		fr := fmt.Sprintf("%d%s%s%s%s%s%s", csvReportItem.EmployeeID, Delimiter,
			csvReportItem.Name, Delimiter, csvReportItem.Department,
			Delimiter, csvReportItem.Manager)

		companyCsvReport = append(companyCsvReport, fr)
	}

	return companyCsvReport
}

// sort based on the department first,
// after that, for each department

func main() {
	solution()
}
